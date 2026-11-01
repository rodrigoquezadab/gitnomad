package gitops

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CloneOptions controla cómo se clona un repo.
type CloneOptions struct {
	// UseGHCredentials hace que git use la sesión de `gh` para autenticarse
	// en github.com solo durante este comando, sin tocar la configuración
	// global. Así los repos privados se clonan aunque no se haya corrido
	// `gh auth setup-git`.
	UseGHCredentials bool
}

// ghCredentialArgs devuelve los -c que configuran gh como helper de
// credenciales para github.com. El primer valor vacío limpia cualquier
// helper previo para ese host y así se evita que otro pida la clave.
func ghCredentialArgs() []string {
	return []string{
		"-c", "credential.https://github.com.helper=",
		"-c", "credential.https://github.com.helper=!gh auth git-credential",
	}
}

// Clone clona url en dest. Si falla, no deja una carpeta a medias.
func Clone(ctx context.Context, url, dest string, opts CloneOptions) error {
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("la carpeta %s ya existe", dest)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("no se pudo crear %s: %w", filepath.Dir(dest), err)
	}

	var args []string
	if opts.UseGHCredentials {
		args = append(args, ghCredentialArgs()...)
	}
	args = append(args, "clone", "--quiet")
	if runtime.GOOS == "windows" {
		// Guardado en el .git/config del repo: rutas de más de 260 caracteres.
		args = append(args, "--config", "core.longpaths=true")
	}
	args = append(args, "--", url, dest)

	if _, err := Run(ctx, "", "git", args...); err != nil {
		_ = os.RemoveAll(dest) // limpiar clon incompleto
		return err
	}
	return nil
}

// IsGitRepo indica si dir es la raíz de un repo git (.git puede ser carpeta
// o archivo, como en worktrees).
func IsGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// OriginURL devuelve la URL del remoto "origin" del repo en dir.
func OriginURL(ctx context.Context, dir string) (string, error) {
	out, err := Run(ctx, dir, "git", "remote", "get-url", "origin")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// NormalizeURL convierte las distintas formas de URL de un repo a una forma
// comparable, por ejemplo:
//
//	https://github.com/Owner/Repo.git  -> github.com/owner/repo
//	git@github.com:Owner/Repo.git      -> github.com/owner/repo
//	ssh://git@github.com/Owner/Repo    -> github.com/owner/repo
func NormalizeURL(u string) string {
	s := strings.TrimSpace(u)
	s = strings.TrimRight(s, "/")
	s = strings.TrimSuffix(s, ".git")
	lower := strings.ToLower(s)

	if i := strings.Index(lower, "://"); i >= 0 {
		scheme := lower[:i]
		rest := s[i+3:]
		if scheme == "file" {
			return "file://" + filepath.ToSlash(filepath.Clean(filepath.FromSlash(rest)))
		}
		if at := strings.Index(rest, "@"); at >= 0 && at < strings.Index(rest+"/", "/") {
			rest = rest[at+1:]
		}
		return strings.ToLower(rest)
	}
	// Forma scp de SSH: usuario@host:ruta
	if at := strings.Index(s, "@"); at >= 0 {
		if colon := strings.Index(s[at:], ":"); colon >= 0 {
			host := s[at+1 : at+colon]
			path := s[at+colon+1:]
			return strings.ToLower(host + "/" + strings.TrimLeft(path, "/"))
		}
	}
	return strings.ToLower(s)
}
