// Package syncer decide qué hacer con cada repo y lo ejecuta en paralelo.
//
// Es el motor compartido: la CLI lo usa hoy y la interfaz gráfica lo usará
// después. Por eso informa el avance con un callback en vez de imprimir.
package syncer

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rodrigoquezadab/gitnomad/internal/config"
	"github.com/rodrigoquezadab/gitnomad/internal/github"
	"github.com/rodrigoquezadab/gitnomad/internal/gitops"
)

// State es el estado local de un repo antes de actuar.
type State string

const (
	// StateMissing: no está en este equipo; hay que clonarlo.
	StateMissing State = "falta"
	// StatePresent: ya está clonado y su origin apunta al repo correcto.
	StatePresent State = "descargado"
	// StateConflict: existe una carpeta con ese nombre pero no es este repo.
	StateConflict State = "conflicto"
)

// Item es un repo con su ruta local y estado.
type Item struct {
	Repo   github.Repo
	Path   string
	State  State
	Reason string // explicación cuando State es conflicto
}

// Skipped es un repo que se dejó fuera por la configuración.
type Skipped struct {
	Repo   github.Repo
	Reason string
}

// Filter separa los repos según la configuración (archivados, forks,
// exclusiones).
func Filter(repos []github.Repo, cfg config.Config) (kept []github.Repo, skipped []Skipped) {
	for _, r := range repos {
		switch {
		case r.IsArchived && !cfg.IncludeArchived:
			skipped = append(skipped, Skipped{r, "archivado"})
		case r.IsFork && !cfg.IncludeForks:
			skipped = append(skipped, Skipped{r, "fork"})
		case cfg.IsExcluded(r.Name):
			skipped = append(skipped, Skipped{r, "excluido"})
		default:
			kept = append(kept, r)
		}
	}
	return kept, skipped
}

// Inspect revisa en disco el estado de cada repo dentro de baseDir.
func Inspect(ctx context.Context, repos []github.Repo, baseDir string) []Item {
	items := make([]Item, len(repos))
	for i, r := range repos {
		items[i] = inspectOne(ctx, r, filepath.Join(baseDir, r.Name))
	}
	return items
}

func inspectOne(ctx context.Context, r github.Repo, path string) Item {
	it := Item{Repo: r, Path: path}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		it.State = StateMissing
		return it
	}
	if err != nil {
		it.State, it.Reason = StateConflict, err.Error()
		return it
	}
	if !info.IsDir() {
		it.State, it.Reason = StateConflict, "existe un archivo con ese nombre"
		return it
	}
	if !gitops.IsGitRepo(path) {
		if empty, _ := isEmptyDir(path); empty {
			// Carpeta vacía: se puede clonar encima sin perder nada.
			it.State = StateMissing
			return it
		}
		it.State, it.Reason = StateConflict, "la carpeta existe pero no es un repo git"
		return it
	}
	origin, err := gitops.OriginURL(ctx, path)
	if err != nil {
		it.State, it.Reason = StateConflict, "el repo no tiene remoto origin"
		return it
	}
	if gitops.NormalizeURL(origin) != gitops.NormalizeURL(r.CloneURL()) {
		it.State, it.Reason = StateConflict, "origin apunta a otro repo: "+origin
		return it
	}
	it.State = StatePresent
	return it
}

func isEmptyDir(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// Result es el resultado de procesar un Item.
type Result struct {
	Item     Item
	Err      error
	Duration time.Duration
}

// Options controla la ejecución.
type Options struct {
	Parallel int
	Clone    gitops.CloneOptions
	// OnStart y OnDone se llaman (desde varias goroutines, pero nunca a la
	// vez) al empezar y terminar cada repo. Pueden ser nil.
	OnStart func(Item)
	OnDone  func(Result)
}

// CloneMissing clona en paralelo los items en estado StateMissing y
// devuelve un resultado por cada uno, en el mismo orden.
func CloneMissing(ctx context.Context, items []Item, opts Options) []Result {
	var todo []Item
	for _, it := range items {
		if it.State == StateMissing {
			todo = append(todo, it)
		}
	}
	results := make([]Result, len(todo))
	if len(todo) == 0 {
		return results
	}
	workers := opts.Parallel
	if workers < 1 {
		workers = 1
	}
	if workers > len(todo) {
		workers = len(todo)
	}

	var mu sync.Mutex // serializa los callbacks
	call := func(f func()) { mu.Lock(); defer mu.Unlock(); f() }

	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				it := todo[i]
				if opts.OnStart != nil {
					call(func() { opts.OnStart(it) })
				}
				start := time.Now()
				err := ctx.Err()
				if err == nil {
					// Si quedó una carpeta vacía, se quita para clonar encima.
					if empty, _ := isEmptyDir(it.Path); empty {
						_ = os.Remove(it.Path)
					}
					err = gitops.Clone(ctx, it.Repo.CloneURL(), it.Path, opts.Clone)
				}
				res := Result{Item: it, Err: err, Duration: time.Since(start)}
				if err == nil {
					res.Item.State = StatePresent
				}
				results[i] = res
				if opts.OnDone != nil {
					call(func() { opts.OnDone(res) })
				}
			}
		}()
	}
	for i := range todo {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results
}
