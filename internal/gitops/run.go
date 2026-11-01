// Package gitops ejecuta git (y gh) de forma segura y multiplataforma.
package gitops

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run ejecuta un comando sin abrir prompts interactivos (para que nunca se
// quede colgado esperando una contraseña) y devuelve stdout. Si falla,
// el error incluye stderr para que el mensaje sea útil.
func Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0", // git no pide usuario/clave por consola
		"GCM_INTERACTIVE=never", // Git Credential Manager (Windows) tampoco
		"GH_PROMPT_DISABLED=1",  // gh no hace preguntas
		"LC_ALL=C",              // salida de git en inglés, estable para parsear
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		if ctx.Err() != nil {
			return stdout.String(), fmt.Errorf("%s %s: cancelado: %w", name, firstArg(args), ctx.Err())
		}
		if msg != "" {
			return stdout.String(), fmt.Errorf("%s %s: %s", name, firstArg(args), msg)
		}
		return stdout.String(), fmt.Errorf("%s %s: %w", name, firstArg(args), err)
	}
	return stdout.String(), nil
}

func firstArg(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && !strings.Contains(a, "=") {
			return a
		}
	}
	return ""
}
