// Package github obtiene la lista de repos del usuario usando GitHub CLI (gh).
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/rodrigoquezadab/gitnomad/internal/gitops"
)

// Repo es un repositorio de GitHub.
type Repo struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	URL           string `json:"url"`
	IsPrivate     bool   `json:"isPrivate"`
	IsFork        bool   `json:"isFork"`
	IsArchived    bool   `json:"isArchived"`
	IsEmpty       bool   `json:"isEmpty"`
	DefaultBranch string `json:"defaultBranch"`
}

// CloneURL devuelve la URL HTTPS para clonar.
func (r Repo) CloneURL() string {
	u := strings.TrimRight(r.URL, "/")
	if strings.HasPrefix(u, "https://") && !strings.HasSuffix(u, ".git") {
		return u + ".git"
	}
	return u
}

// Lister obtiene repos. Es una interfaz para poder probar el resto sin red.
type Lister interface {
	ListRepos(ctx context.Context) ([]Repo, error)
}

// GH lista repos con `gh repo list`, que usa la sesión iniciada con
// `gh auth login` e incluye repos públicos y privados del usuario.
type GH struct {
	// Limit es el máximo de repos a pedir (por defecto 5000).
	Limit int
}

type ghRepo struct {
	Name             string `json:"name"`
	NameWithOwner    string `json:"nameWithOwner"`
	URL              string `json:"url"`
	IsPrivate        bool   `json:"isPrivate"`
	IsFork           bool   `json:"isFork"`
	IsArchived       bool   `json:"isArchived"`
	IsEmpty          bool   `json:"isEmpty"`
	DefaultBranchRef *struct {
		Name string `json:"name"`
	} `json:"defaultBranchRef"`
}

// ListRepos devuelve los repos del usuario autenticado ordenados por nombre.
func (g GH) ListRepos(ctx context.Context) ([]Repo, error) {
	limit := g.Limit
	if limit <= 0 {
		limit = 5000
	}
	out, err := gitops.Run(ctx, "", "gh", "repo", "list",
		"--limit", fmt.Sprint(limit),
		"--json", "name,nameWithOwner,url,isPrivate,isFork,isArchived,isEmpty,defaultBranchRef")
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la lista de repos de GitHub "+
			"(¿iniciaste sesión con `gh auth login`?): %w", err)
	}
	return parse([]byte(out))
}

func parse(data []byte) ([]Repo, error) {
	var raw []ghRepo
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("respuesta inesperada de gh: %w", err)
	}
	repos := make([]Repo, 0, len(raw))
	for _, r := range raw {
		repo := Repo{
			Name:          r.Name,
			NameWithOwner: r.NameWithOwner,
			URL:           r.URL,
			IsPrivate:     r.IsPrivate,
			IsFork:        r.IsFork,
			IsArchived:    r.IsArchived,
			IsEmpty:       r.IsEmpty,
		}
		if r.DefaultBranchRef != nil {
			repo.DefaultBranch = r.DefaultBranchRef.Name
		}
		repos = append(repos, repo)
	}
	sort.Slice(repos, func(i, j int) bool {
		return strings.ToLower(repos[i].Name) < strings.ToLower(repos[j].Name)
	})
	return repos, nil
}

// CurrentUser devuelve el usuario con sesión iniciada en gh.
func CurrentUser(ctx context.Context) (string, error) {
	out, err := gitops.Run(ctx, "", "gh", "api", "user", "--jq", ".login")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Static es un Lister con una lista fija (útil para pruebas).
type Static []Repo

// ListRepos devuelve la lista fija.
func (s Static) ListRepos(context.Context) ([]Repo, error) { return []Repo(s), nil }
