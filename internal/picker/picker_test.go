package picker

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	none := func(string) (string, error) { return "", exec.ErrNotFound }
	only := func(tool string) func(string) (string, error) {
		return func(name string) (string, error) {
			if name == tool {
				return "/usr/bin/" + name, nil
			}
			return "", exec.ErrNotFound
		}
	}

	name, _, env, err := command("windows", "Elige", `C:\dev`, none)
	if err != nil || name != "powershell.exe" || !strings.Contains(strings.Join(env, "\n"), `GITNOMAD_PICK_START=C:\dev`) {
		t.Errorf("windows: %s %v %v", name, env, err)
	}

	name, args, _, err := command("darwin", "Elige", "", none)
	if err != nil || name != "osascript" || strings.Contains(strings.Join(args, " "), "default location") {
		t.Errorf("darwin sin start: %s %v %v", name, args, err)
	}

	name, args, _, err = command("linux", "Elige", "/home/a", only("kdialog"))
	if err != nil || name != "kdialog" || args[len(args)-1] != "/home/a" {
		t.Errorf("linux kdialog: %s %v %v", name, args, err)
	}

	if name, _, _, _ = command("linux", "Elige", "", only("zenity")); name != "zenity" {
		t.Errorf("linux zenity: %s", name)
	}

	if _, _, _, err = command("linux", "Elige", "", none); !errors.Is(err, ErrUnavailable) {
		t.Errorf("linux sin herramientas: %v", err)
	}
}

func TestResult(t *testing.T) {
	if got, err := result("  /home/a/dev/ \n", "", nil); err != nil || got != filepath.Clean("/home/a/dev/") {
		t.Errorf("ruta: %q %v", got, err)
	}
	if _, err := result("", "", nil); !errors.Is(err, ErrCanceled) {
		t.Errorf("vacío: %v", err)
	}
	if _, err := result("", "", errors.New("boom")); err == nil || errors.Is(err, ErrCanceled) {
		t.Errorf("error: %v", err)
	}
}

func TestNearestDir(t *testing.T) {
	dir := t.TempDir()
	if got := nearestDir(filepath.Join(dir, "no", "existe")); got != dir {
		t.Errorf("nearestDir: %q, se esperaba %q", got, dir)
	}
	if got := nearestDir(""); got != "" {
		t.Errorf("vacío: %q", got)
	}
}
