// Package picker abre el diálogo nativo del sistema para elegir una carpeta.
//
// No usa dependencias: en cada sistema llama a una herramienta que ya viene
// instalada (o casi siempre lo está):
//   - Windows: PowerShell con System.Windows.Forms.FolderBrowserDialog
//   - macOS:   osascript ("choose folder")
//   - Linux:   zenity o, si no está, kdialog
package picker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrCanceled indica que el usuario cerró el diálogo sin elegir carpeta.
var ErrCanceled = errors.New("no se eligió ninguna carpeta")

// ErrUnavailable indica que no hay forma de mostrar un diálogo en este equipo
// (por ejemplo, un servidor Linux sin escritorio).
var ErrUnavailable = errors.New("no se puede abrir un selector de carpetas en este equipo")

const winScript = `Add-Type -AssemblyName System.Windows.Forms
[Console]::OutputEncoding = [Text.Encoding]::UTF8
$d = New-Object System.Windows.Forms.FolderBrowserDialog
$d.Description = $env:GITNOMAD_PICK_TITLE
$d.ShowNewFolderButton = $true
if ($env:GITNOMAD_PICK_START) { $d.SelectedPath = $env:GITNOMAD_PICK_START }
$owner = New-Object System.Windows.Forms.Form -Property @{TopMost = $true}
if ($d.ShowDialog($owner) -eq 'OK') { $d.SelectedPath }`

// Folder muestra el diálogo y devuelve la ruta absoluta elegida. start es la
// carpeta donde abre el diálogo; si no existe se usa la más cercana que sí.
func Folder(ctx context.Context, title, start string) (string, error) {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" &&
		os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return "", fmt.Errorf("%w (no hay escritorio gráfico)", ErrUnavailable)
	}
	name, args, env, err := command(runtime.GOOS, title, nearestDir(start), exec.LookPath)
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = append(os.Environ(), env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	return result(stdout.String(), stderr.String(), runErr)
}

// command arma la llamada para cada sistema. Está separada para poder
// probarla sin abrir ventanas.
func command(goos, title, start string, lookPath func(string) (string, error)) (name string, args, env []string, err error) {
	switch goos {
	case "windows":
		return "powershell.exe", []string{"-NoProfile", "-NonInteractive", "-STA", "-Command", winScript},
			[]string{"GITNOMAD_PICK_TITLE=" + title, "GITNOMAD_PICK_START=" + start}, nil
	case "darwin":
		script := []string{"-e", "on run argv"}
		if start != "" {
			script = append(script, "-e", "POSIX path of (choose folder with prompt (item 1 of argv) default location (POSIX file (item 2 of argv)))")
		} else {
			script = append(script, "-e", "POSIX path of (choose folder with prompt (item 1 of argv))")
		}
		script = append(script, "-e", "end run", title, start)
		return "osascript", script, nil, nil
	default:
		if _, e := lookPath("zenity"); e == nil {
			args := []string{"--file-selection", "--directory", "--title=" + title}
			if start != "" {
				args = append(args, "--filename="+start+string(filepath.Separator))
			}
			return "zenity", args, nil, nil
		}
		if _, e := lookPath("kdialog"); e == nil {
			return "kdialog", []string{"--title", title, "--getexistingdirectory", start}, nil, nil
		}
		return "", nil, nil, fmt.Errorf("%w (instala zenity o kdialog)", ErrUnavailable)
	}
}

// result interpreta la salida del diálogo. Todas las herramientas salen con
// código 1 (y osascript con el error -128) cuando el usuario cancela.
func result(stdout, stderr string, runErr error) (string, error) {
	path := strings.TrimSpace(stdout)
	if runErr != nil {
		var ee *exec.ExitError
		if errors.As(runErr, &ee) && ee.ExitCode() == 1 && (strings.TrimSpace(stderr) == "" || strings.Contains(stderr, "-128")) {
			return "", ErrCanceled
		}
		if msg := strings.TrimSpace(stderr); msg != "" {
			return "", fmt.Errorf("selector de carpetas: %s", msg)
		}
		return "", fmt.Errorf("selector de carpetas: %w", runErr)
	}
	if path == "" {
		return "", ErrCanceled
	}
	return filepath.Clean(path), nil
}

// nearestDir devuelve dir si existe o, si no, la carpeta padre más cercana
// que exista. Devuelve "" si no encuentra ninguna.
func nearestDir(dir string) string {
	if dir == "" {
		return ""
	}
	for d := filepath.Clean(dir); ; d = filepath.Dir(d) {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			return d
		}
		if filepath.Dir(d) == d {
			return ""
		}
	}
}
