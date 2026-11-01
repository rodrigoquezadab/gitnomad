// Package doctor revisa que el equipo tenga todo lo necesario y, si se pide,
// arregla lo que se pueda.
package doctor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rodrigoquezadab/gitnomad/internal/config"
	"github.com/rodrigoquezadab/gitnomad/internal/github"
	"github.com/rodrigoquezadab/gitnomad/internal/gitops"
)

// Level es la gravedad de un chequeo.
type Level int

const (
	OK Level = iota
	Warn
	Fail
)

// Check es el resultado de un chequeo.
type Check struct {
	Name   string
	Level  Level
	Detail string
	Hint   string // qué hacer para arreglarlo
	Fixed  bool   // se arregló durante esta ejecución
}

// Run ejecuta todos los chequeos. Con fix=true intenta arreglar los
// problemas que se pueden resolver sin intervención.
func Run(ctx context.Context, cfg config.Config, fix bool) []Check {
	var checks []Check

	gitOK := true
	if v, err := version(ctx, "git", "--version"); err != nil {
		gitOK = false
		checks = append(checks, Check{Name: "git instalado", Level: Fail, Detail: err.Error(),
			Hint: installHint("git")})
	} else {
		checks = append(checks, Check{Name: "git instalado", Level: OK, Detail: v})
	}

	ghOK := true
	if v, err := version(ctx, "gh", "--version"); err != nil {
		ghOK = false
		checks = append(checks, Check{Name: "GitHub CLI (gh) instalado", Level: Fail, Detail: err.Error(),
			Hint: installHint("gh")})
	} else {
		checks = append(checks, Check{Name: "GitHub CLI (gh) instalado", Level: OK, Detail: firstLine(v)})
	}

	loggedIn := false
	if ghOK {
		if user, err := github.CurrentUser(ctx); err != nil {
			checks = append(checks, Check{Name: "Sesión en GitHub", Level: Fail,
				Detail: "no hay sesión iniciada en gh",
				Hint:   "corre: gh auth login   (elige GitHub.com y HTTPS)"})
		} else {
			loggedIn = true
			checks = append(checks, Check{Name: "Sesión en GitHub", Level: OK, Detail: "usuario " + user})
		}
	}

	if gitOK && ghOK && loggedIn {
		checks = append(checks, credentialCheck(ctx, fix))
	}

	if gitOK && runtime.GOOS == "windows" {
		checks = append(checks, longPathsCheck(ctx, fix))
	}

	checks = append(checks, baseDirCheck(cfg.BaseDir, fix))
	return checks
}

// credentialCheck: para que `git push/pull` en los repos clonados use la
// sesión de gh también fuera de gitnomad.
func credentialCheck(ctx context.Context, fix bool) Check {
	c := Check{Name: "git usa las credenciales de gh"}
	out, _ := gitops.Run(ctx, "", "git", "config", "--global", "--get-all", "credential.https://github.com.helper")
	if strings.Contains(out, "gh") && strings.Contains(out, "auth git-credential") {
		c.Level, c.Detail = OK, "configurado"
		return c
	}
	if fix {
		if _, err := gitops.Run(ctx, "", "gh", "auth", "setup-git"); err == nil {
			c.Level, c.Detail, c.Fixed = OK, "configurado con `gh auth setup-git`", true
			return c
		} else {
			c.Detail = err.Error()
		}
	}
	c.Level = Warn
	if c.Detail == "" {
		c.Detail = "no configurado (gitnomad igual puede clonar, pero `git push` manual podría pedir clave)"
	}
	c.Hint = "corre: gitnomad doctor --fix   (o: gh auth setup-git)"
	return c
}

func longPathsCheck(ctx context.Context, fix bool) Check {
	c := Check{Name: "Rutas largas en Windows (core.longpaths)"}
	out, _ := gitops.Run(ctx, "", "git", "config", "--global", "--get", "core.longpaths")
	if strings.TrimSpace(out) == "true" {
		c.Level, c.Detail = OK, "activado"
		return c
	}
	if fix {
		if _, err := gitops.Run(ctx, "", "git", "config", "--global", "core.longpaths", "true"); err == nil {
			c.Level, c.Detail, c.Fixed = OK, "activado", true
			return c
		}
	}
	c.Level, c.Detail = Warn, "desactivado (los repos clonados por gitnomad lo activan igual)"
	c.Hint = "corre: gitnomad doctor --fix"
	return c
}

func baseDirCheck(dir string, fix bool) Check {
	c := Check{Name: "Carpeta de repos"}
	info, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		if fix {
			if err := os.MkdirAll(dir, 0o755); err == nil {
				c.Level, c.Detail, c.Fixed = OK, dir+" (creada)", true
				return c
			}
		}
		c.Level, c.Detail = Warn, dir+" no existe (se creará al clonar)"
		return c
	case err != nil:
		c.Level, c.Detail = Fail, err.Error()
		return c
	case !info.IsDir():
		c.Level, c.Detail = Fail, dir+" existe pero no es una carpeta"
		c.Hint = "cambia la carpeta con: gitnomad config set base_dir <ruta>"
		return c
	}
	probe := filepath.Join(dir, ".gitnomad-write-test")
	if err := os.WriteFile(probe, []byte("ok"), 0o644); err != nil {
		c.Level, c.Detail = Fail, "sin permiso de escritura en "+dir
		return c
	}
	_ = os.Remove(probe)
	c.Level, c.Detail = OK, dir
	return c
}

func version(ctx context.Context, name string, args ...string) (string, error) {
	if _, err := exec.LookPath(name); err != nil {
		return "", fmt.Errorf("%s no está instalado o no está en el PATH", name)
	}
	out, err := gitops.Run(ctx, "", name, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func installHint(tool string) string {
	switch runtime.GOOS + "/" + tool {
	case "windows/git":
		return "instala con: winget install --id Git.Git -e"
	case "windows/gh":
		return "instala con: winget install --id GitHub.cli -e"
	case "darwin/git":
		return "instala con: xcode-select --install   (o: brew install git)"
	case "darwin/gh":
		return "instala con: brew install gh"
	case "linux/git":
		return "instala con el gestor de paquetes, ej: sudo apt install git"
	default:
		return "instala con el gestor de paquetes; ver https://github.com/cli/cli#installation"
	}
}
