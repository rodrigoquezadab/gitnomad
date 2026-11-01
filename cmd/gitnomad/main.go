// gitnomad descarga y mantiene al día todos tus repos de GitHub en cualquier
// equipo (Windows, macOS o Linux).
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rodrigoquezadab/gitnomad/internal/config"
	"github.com/rodrigoquezadab/gitnomad/internal/doctor"
	"github.com/rodrigoquezadab/gitnomad/internal/github"
	"github.com/rodrigoquezadab/gitnomad/internal/gitops"
	"github.com/rodrigoquezadab/gitnomad/internal/syncer"
)

// version se reemplaza al compilar: -ldflags "-X main.version=v0.1.0"
var version = "dev"

const usage = `gitnomad - descarga todos tus repos de GitHub en este equipo

Uso:
  gitnomad <comando> [opciones]

Comandos:
  doctor          Revisa que git, gh y la sesión de GitHub estén listos
                  (--fix arregla lo que se pueda)
  config          Muestra o cambia la configuración de este equipo
                    gitnomad config show
                    gitnomad config init [--dir RUTA]
                    gitnomad config set <clave> <valor>
  list            Lista tus repos y si ya están en este equipo (--json)
  clone           Clona los repos que faltan en este equipo
                    --dry-run      muestra qué haría sin clonar
                    --parallel N   repos a la vez
                    --dir RUTA     usa otra carpeta solo esta vez
  version         Muestra la versión

Primeros pasos en un equipo nuevo:
  gh auth login
  gitnomad doctor --fix
  gitnomad config init --dir ~/dev/github
  gitnomad clone
`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	os.Exit(run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, out, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(out, usage)
		return 0
	}
	cmd, rest := args[0], args[1:]
	var err error
	switch cmd {
	case "doctor":
		err = cmdDoctor(ctx, rest, out)
	case "config":
		err = cmdConfig(rest, out)
	case "list", "ls":
		err = cmdList(ctx, rest, out, github.GH{})
	case "clone":
		err = cmdClone(ctx, rest, out, github.GH{})
	case "version", "--version", "-v":
		fmt.Fprintf(out, "gitnomad %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
	case "help", "--help", "-h":
		fmt.Fprint(out, usage)
	default:
		fmt.Fprintf(errOut, "comando desconocido: %q\n\n%s", cmd, usage)
		return 2
	}
	if err != nil {
		var ue usageError
		if errors.As(err, &ue) {
			fmt.Fprintf(errOut, "error: %v\n", err)
			return 2
		}
		if !errors.Is(err, errReported) {
			fmt.Fprintf(errOut, "error: %v\n", err)
		}
		return 1
	}
	return 0
}

type usageError struct{ msg string }

func (e usageError) Error() string { return e.msg }

// errReported indica que el detalle ya se mostró y solo hay que salir con 1.
var errReported = errors.New("errores reportados")

func newFlags(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

func parseFlags(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		return usageError{err.Error()}
	}
	return nil
}

func loadConfig(out io.Writer) (config.Config, error) {
	cfg, exists, err := config.Load()
	if err != nil {
		return cfg, err
	}
	if !exists {
		fmt.Fprintf(out, "(sin configuración en este equipo; usando la carpeta %s. "+
			"Cámbiala con `gitnomad config init --dir RUTA`)\n\n", cfg.BaseDir)
	}
	return cfg, nil
}

// ---------- doctor ----------

func cmdDoctor(ctx context.Context, args []string, out io.Writer) error {
	fs := newFlags("doctor")
	fix := fs.Bool("fix", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	cfg, _, err := config.Load()
	if err != nil {
		fmt.Fprintf(out, "aviso: %v\n", err)
	}
	failed := false
	for _, c := range doctor.Run(ctx, cfg, *fix) {
		tag := map[doctor.Level]string{doctor.OK: "[ OK ]", doctor.Warn: "[AVISO]", doctor.Fail: "[FALLA]"}[c.Level]
		line := fmt.Sprintf("%-7s %s: %s", tag, c.Name, c.Detail)
		if c.Fixed {
			line += "  <- arreglado"
		}
		fmt.Fprintln(out, line)
		if c.Hint != "" && c.Level != doctor.OK {
			fmt.Fprintf(out, "        %s\n", c.Hint)
		}
		if c.Level == doctor.Fail {
			failed = true
		}
	}
	if failed {
		fmt.Fprintln(out, "\nHay problemas que resolver antes de usar gitnomad.")
		return errReported
	}
	fmt.Fprintln(out, "\nTodo listo.")
	return nil
}

// ---------- config ----------

func cmdConfig(args []string, out io.Writer) error {
	sub := "show"
	if len(args) > 0 {
		sub, args = args[0], args[1:]
	}
	switch sub {
	case "show":
		cfg, exists, err := config.Load()
		if err != nil {
			return err
		}
		p, _ := config.Path()
		if exists {
			fmt.Fprintf(out, "Archivo: %s\n\n", p)
		} else {
			fmt.Fprintf(out, "Archivo: %s (no existe; valores por defecto)\n\n", p)
		}
		printConfig(out, cfg)
		return nil

	case "init":
		fs := newFlags("config init")
		dir := fs.String("dir", "", "")
		force := fs.Bool("force", false, "")
		if err := parseFlags(fs, args); err != nil {
			return err
		}
		_, exists, _ := config.Load()
		if exists && !*force {
			return usageError{"ya existe una configuración; usa `gitnomad config set` o agrega --force"}
		}
		cfg := config.Default()
		if *dir != "" {
			cfg.BaseDir = *dir
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		p, _ := config.Path()
		fmt.Fprintf(out, "Configuración guardada en %s\n\n", p)
		cfg, _, _ = config.Load()
		printConfig(out, cfg)
		return nil

	case "set":
		if len(args) != 2 {
			return usageError{"uso: gitnomad config set <clave> <valor>\n" +
				"claves: base_dir, parallel, include_archived, include_forks, exclude (separado por comas)"}
		}
		cfg, _, err := config.Load()
		if err != nil {
			return err
		}
		key, val := args[0], args[1]
		switch key {
		case "base_dir":
			cfg.BaseDir = val
		case "parallel":
			n, err := strconv.Atoi(val)
			if err != nil {
				return usageError{"parallel debe ser un número"}
			}
			cfg.Parallel = n
		case "include_archived", "include_forks":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return usageError{key + " debe ser true o false"}
			}
			if key == "include_archived" {
				cfg.IncludeArchived = b
			} else {
				cfg.IncludeForks = b
			}
		case "exclude":
			cfg.Exclude = strings.Split(val, ",")
		default:
			return usageError{"clave desconocida: " + key}
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Fprintln(out, "Guardado.")
		cfg, _, _ = config.Load()
		printConfig(out, cfg)
		return nil
	}
	return usageError{"subcomando de config desconocido: " + sub}
}

func printConfig(out io.Writer, c config.Config) {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "base_dir\t%s\n", c.BaseDir)
	fmt.Fprintf(tw, "parallel\t%d\n", c.Parallel)
	fmt.Fprintf(tw, "include_archived\t%t\n", c.IncludeArchived)
	fmt.Fprintf(tw, "include_forks\t%t\n", c.IncludeForks)
	fmt.Fprintf(tw, "exclude\t%s\n", strings.Join(c.Exclude, ","))
	tw.Flush()
}

// ---------- list ----------

type listEntry struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Fork     bool   `json:"fork"`
	Archived bool   `json:"archived"`
	State    string `json:"state"`
	Path     string `json:"path,omitempty"`
	Note     string `json:"note,omitempty"`
}

func cmdList(ctx context.Context, args []string, out io.Writer, lister github.Lister) error {
	fs := newFlags("list")
	asJSON := fs.Bool("json", false, "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	var cfg config.Config
	var err error
	if *asJSON {
		cfg, _, err = config.Load()
	} else {
		cfg, err = loadConfig(out)
	}
	if err != nil {
		return err
	}
	repos, err := lister.ListRepos(ctx)
	if err != nil {
		return err
	}
	kept, skipped := syncer.Filter(repos, cfg)
	items := syncer.Inspect(ctx, kept, cfg.BaseDir)

	var entries []listEntry
	for _, it := range items {
		entries = append(entries, listEntry{Name: it.Repo.Name, Private: it.Repo.IsPrivate,
			Fork: it.Repo.IsFork, Archived: it.Repo.IsArchived, State: string(it.State),
			Path: it.Path, Note: it.Reason})
	}
	for _, s := range skipped {
		entries = append(entries, listEntry{Name: s.Repo.Name, Private: s.Repo.IsPrivate,
			Fork: s.Repo.IsFork, Archived: s.Repo.IsArchived, State: "omitido", Note: s.Reason})
	}

	if *asJSON {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if entries == nil {
			entries = []listEntry{}
		}
		return enc.Encode(entries)
	}

	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "REPO\tVISIBILIDAD\tESTADO\tNOTA")
	counts := map[string]int{}
	for _, e := range entries {
		vis := "público"
		if e.Private {
			vis = "privado"
		}
		var tags []string
		if e.Fork {
			tags = append(tags, "fork")
		}
		if e.Archived {
			tags = append(tags, "archivado")
		}
		note := e.Note
		if len(tags) > 0 && e.State != "omitido" {
			note = strings.TrimSpace(strings.Join(tags, ", ") + " " + note)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", e.Name, vis, e.State, note)
		counts[e.State]++
	}
	tw.Flush()
	fmt.Fprintf(out, "\n%d repos: %d descargados, %d faltan, %d en conflicto, %d omitidos\n",
		len(entries), counts[string(syncer.StatePresent)], counts[string(syncer.StateMissing)],
		counts[string(syncer.StateConflict)], counts["omitido"])
	fmt.Fprintf(out, "Carpeta: %s\n", cfg.BaseDir)
	return nil
}

// ---------- clone ----------

func cmdClone(ctx context.Context, args []string, out io.Writer, lister github.Lister) error {
	fs := newFlags("clone")
	dry := fs.Bool("dry-run", false, "")
	parallel := fs.Int("parallel", 0, "")
	dir := fs.String("dir", "", "")
	if err := parseFlags(fs, args); err != nil {
		return err
	}
	cfg, err := loadConfig(out)
	if err != nil {
		return err
	}
	if *dir != "" {
		cfg.BaseDir = config.ExpandHome(*dir)
	}
	if *parallel > 0 {
		cfg.Parallel = *parallel
	}
	if err := cfg.Validate(); err != nil {
		return usageError{err.Error()}
	}

	fmt.Fprintln(out, "Obteniendo tus repos de GitHub...")
	repos, err := lister.ListRepos(ctx)
	if err != nil {
		return err
	}
	kept, skipped := syncer.Filter(repos, cfg)
	items := syncer.Inspect(ctx, kept, cfg.BaseDir)

	var missing, present, conflicts []syncer.Item
	for _, it := range items {
		switch it.State {
		case syncer.StateMissing:
			missing = append(missing, it)
		case syncer.StatePresent:
			present = append(present, it)
		case syncer.StateConflict:
			conflicts = append(conflicts, it)
		}
	}
	fmt.Fprintf(out, "%d repos en GitHub: %d ya descargados, %d por clonar, %d en conflicto, %d omitidos\n",
		len(repos), len(present), len(missing), len(conflicts), len(skipped))
	fmt.Fprintf(out, "Carpeta: %s\n\n", cfg.BaseDir)

	for _, c := range conflicts {
		fmt.Fprintf(out, "  conflicto  %s: %s\n", c.Repo.Name, c.Reason)
	}
	if len(conflicts) > 0 {
		fmt.Fprintln(out, "  (los repos en conflicto no se tocan; revísalos a mano)")
		fmt.Fprintln(out)
	}

	if len(missing) == 0 {
		fmt.Fprintln(out, "No hay nada que clonar. Todo está descargado.")
		return nil
	}
	if *dry {
		fmt.Fprintln(out, "Se clonarían (--dry-run, no se hizo nada):")
		for _, m := range missing {
			fmt.Fprintf(out, "  %s -> %s\n", m.Repo.NameWithOwner, m.Path)
		}
		return nil
	}

	total := len(missing)
	done := 0
	start := time.Now()
	results := syncer.CloneMissing(ctx, items, syncer.Options{
		Parallel: cfg.Parallel,
		Clone:    gitops.CloneOptions{UseGHCredentials: true},
		OnDone: func(r syncer.Result) {
			done++
			if r.Err != nil {
				fmt.Fprintf(out, "[%d/%d] ERROR  %s: %v\n", done, total, r.Item.Repo.Name, r.Err)
				return
			}
			fmt.Fprintf(out, "[%d/%d] ok     %s (%.1fs)\n", done, total, r.Item.Repo.Name, r.Duration.Seconds())
		},
	})

	failed := 0
	for _, r := range results {
		if r.Err != nil {
			failed++
		}
	}
	fmt.Fprintf(out, "\nListo en %.0fs: %d clonados, %d con error.\n",
		time.Since(start).Seconds(), total-failed, failed)
	if failed > 0 {
		return errReported
	}
	return nil
}
