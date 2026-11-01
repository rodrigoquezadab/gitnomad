package syncer

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rodrigoquezadab/gitnomad/internal/config"
	"github.com/rodrigoquezadab/gitnomad/internal/github"
	"github.com/rodrigoquezadab/gitnomad/internal/gitops"
)

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// makeRemote crea un repo "remoto" bare con un commit y devuelve su URL file://.
func makeRemote(t *testing.T, root, name string) string {
	t.Helper()
	bare := filepath.Join(root, "remotes", name+".git")
	work := filepath.Join(root, "work", name)
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "--bare", "--quiet", "--initial-branch=main", bare)
	git(t, work, "init", "--quiet")
	if err := os.WriteFile(filepath.Join(work, "README.md"), []byte("# "+name), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, work, "add", ".")
	git(t, work, "commit", "--quiet", "-m", "init")
	git(t, work, "push", "--quiet", bare, "HEAD:refs/heads/main")
	return fileURL(bare)
}

// fileURL arma una URL file:// válida también en Windows (file:///C:/...).
func fileURL(p string) string {
	s := filepath.ToSlash(p)
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	return "file://" + s
}

func TestFilter(t *testing.T) {
	repos := []github.Repo{
		{Name: "a"}, {Name: "arch", IsArchived: true}, {Name: "fk", IsFork: true}, {Name: "test-x"},
	}
	cfg := config.Default()
	cfg.Exclude = []string{"test-*"}
	kept, skipped := Filter(repos, cfg)
	if len(kept) != 2 || len(skipped) != 2 {
		t.Fatalf("kept=%v skipped=%v", kept, skipped)
	}
	cfg.IncludeArchived, cfg.IncludeForks, cfg.Exclude = true, false, nil
	kept, _ = Filter(repos, cfg)
	if len(kept) != 3 {
		t.Fatalf("kept=%v", kept)
	}
}

func TestInspectAndClone(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	base := filepath.Join(root, "repos")

	repos := []github.Repo{
		{Name: "alpha", NameWithOwner: "me/alpha", URL: makeRemote(t, root, "alpha")},
		{Name: "beta", NameWithOwner: "me/beta", URL: makeRemote(t, root, "beta")},
		{Name: "gamma", NameWithOwner: "me/gamma", URL: makeRemote(t, root, "gamma")},
		{Name: "otra", NameWithOwner: "me/otra", URL: makeRemote(t, root, "otra")},
		{Name: "vacia", NameWithOwner: "me/vacia", URL: makeRemote(t, root, "vacia")},
		{Name: "malo", NameWithOwner: "me/malo", URL: fileURL(filepath.Join(root, "no-existe.git"))},
	}

	// Preparar casos en disco:
	// gamma: carpeta con archivos que no es repo -> conflicto
	must(t, os.MkdirAll(filepath.Join(base, "gamma"), 0o755))
	must(t, os.WriteFile(filepath.Join(base, "gamma", "notas.txt"), []byte("mío"), 0o644))
	// otra: repo clonado desde otro remoto -> conflicto
	git(t, root, "clone", "--quiet", repos[0].URL, filepath.Join(base, "otra"))
	// vacia: carpeta vacía -> se puede clonar
	must(t, os.MkdirAll(filepath.Join(base, "vacia"), 0o755))

	items := Inspect(ctx, repos, base)
	want := map[string]State{"alpha": StateMissing, "beta": StateMissing, "gamma": StateConflict,
		"otra": StateConflict, "vacia": StateMissing, "malo": StateMissing}
	for _, it := range items {
		if it.State != want[it.Repo.Name] {
			t.Errorf("%s: estado %q, se esperaba %q (%s)", it.Repo.Name, it.State, want[it.Repo.Name], it.Reason)
		}
	}

	var started, finished int
	results := CloneMissing(ctx, items, Options{
		Parallel: 3,
		OnStart:  func(Item) { started++ },
		OnDone:   func(Result) { finished++ },
	})
	if len(results) != 4 || started != 4 || finished != 4 {
		t.Fatalf("results=%d started=%d finished=%d", len(results), started, finished)
	}
	for _, r := range results {
		if r.Item.Repo.Name == "malo" {
			if r.Err == nil {
				t.Error("malo: se esperaba error")
			}
			if _, err := os.Stat(r.Item.Path); !os.IsNotExist(err) {
				t.Error("malo: quedó una carpeta a medias")
			}
			continue
		}
		if r.Err != nil {
			t.Errorf("%s: %v", r.Item.Repo.Name, r.Err)
		}
	}

	// Los archivos del usuario en gamma siguen intactos.
	if b, _ := os.ReadFile(filepath.Join(base, "gamma", "notas.txt")); string(b) != "mío" {
		t.Error("gamma fue modificado")
	}

	// Segunda pasada: todo lo clonado aparece como descargado.
	for _, it := range Inspect(ctx, repos, base) {
		exp := map[string]State{"alpha": StatePresent, "beta": StatePresent, "vacia": StatePresent,
			"gamma": StateConflict, "otra": StateConflict, "malo": StateMissing}[it.Repo.Name]
		if it.State != exp {
			t.Errorf("2da pasada %s: %q, se esperaba %q", it.Repo.Name, it.State, exp)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(base, "alpha", "README.md")); !strings.Contains(string(b), "alpha") {
		t.Error("alpha no tiene el contenido esperado")
	}
}

func TestCloneCancelled(t *testing.T) {
	root := t.TempDir()
	url := makeRemote(t, root, "x")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	items := []Item{{Repo: github.Repo{Name: "x", URL: url}, Path: filepath.Join(root, "out", "x"), State: StateMissing}}
	res := CloneMissing(ctx, items, Options{Parallel: 2, Clone: gitops.CloneOptions{}})
	if len(res) != 1 || res[0].Err == nil {
		t.Fatalf("se esperaba error por cancelación: %+v", res)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
