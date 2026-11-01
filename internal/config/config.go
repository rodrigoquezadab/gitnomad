// Package config carga y guarda la configuración de gitnomad de cada equipo.
//
// El archivo vive en la carpeta de configuración estándar del sistema:
//   - Windows: %APPDATA%\gitnomad\config.json
//   - macOS:   ~/Library/Application Support/gitnomad/config.json
//   - Linux:   ~/.config/gitnomad/config.json  (o $XDG_CONFIG_HOME)
//
// Se puede forzar otra ruta con la variable de entorno GITNOMAD_CONFIG.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Config es la configuración de un equipo. Cada dispositivo tiene la suya,
// así cada uno puede guardar los repos en una carpeta distinta.
type Config struct {
	// BaseDir es la carpeta donde se clonan los repos (uno por subcarpeta).
	BaseDir string `json:"base_dir"`
	// Parallel es cuántos repos se procesan a la vez.
	Parallel int `json:"parallel"`
	// IncludeArchived incluye los repos archivados en GitHub.
	IncludeArchived bool `json:"include_archived"`
	// IncludeForks incluye los repos que son fork de otro.
	IncludeForks bool `json:"include_forks"`
	// Exclude son patrones (estilo glob, sin distinguir mayúsculas) de
	// nombres de repo que se ignoran. Ej: "test-*", "viejo-proyecto".
	Exclude []string `json:"exclude"`
}

// Default devuelve la configuración por defecto: ~/dev/github, 4 en paralelo,
// forks incluidos y archivados excluidos.
func Default() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return Config{
		BaseDir:         filepath.Join(home, "dev", "github"),
		Parallel:        4,
		IncludeArchived: false,
		IncludeForks:    true,
		Exclude:         []string{},
	}
}

// Path devuelve la ruta del archivo de configuración.
func Path() (string, error) {
	if p := os.Getenv("GITNOMAD_CONFIG"); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("no se pudo determinar la carpeta de configuración: %w", err)
	}
	return filepath.Join(dir, "gitnomad", "config.json"), nil
}

// Load lee la configuración. Si el archivo no existe devuelve Default() y
// exists=false (no es un error).
func Load() (cfg Config, exists bool, err error) {
	p, err := Path()
	if err != nil {
		return Default(), false, err
	}
	data, err := os.ReadFile(p)
	if errors.Is(err, fs.ErrNotExist) {
		return Default(), false, nil
	}
	if err != nil {
		return Default(), false, fmt.Errorf("no se pudo leer %s: %w", p, err)
	}
	cfg = Default()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), true, fmt.Errorf("el archivo %s no es JSON válido: %w", p, err)
	}
	cfg.normalize()
	return cfg, true, cfg.Validate()
}

// Save escribe la configuración de forma atómica (archivo temporal + rename).
func Save(cfg Config) error {
	cfg.normalize()
	if err := cfg.Validate(); err != nil {
		return err
	}
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("no se pudo crear %s: %w", filepath.Dir(p), err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("no se pudo escribir %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, p); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("no se pudo guardar %s: %w", p, err)
	}
	return nil
}

// Validate revisa que los valores tengan sentido.
func (c Config) Validate() error {
	if strings.TrimSpace(c.BaseDir) == "" {
		return errors.New("base_dir no puede estar vacío")
	}
	if !filepath.IsAbs(c.BaseDir) {
		return fmt.Errorf("base_dir debe ser una ruta absoluta (recibido: %q)", c.BaseDir)
	}
	if c.Parallel < 1 || c.Parallel > 32 {
		return fmt.Errorf("parallel debe estar entre 1 y 32 (recibido: %d)", c.Parallel)
	}
	for _, pat := range c.Exclude {
		if _, err := filepath.Match(strings.ToLower(pat), ""); err != nil {
			return fmt.Errorf("patrón de exclusión inválido %q: %w", pat, err)
		}
	}
	return nil
}

// IsExcluded indica si un nombre de repo coincide con algún patrón de Exclude.
func (c Config) IsExcluded(name string) bool {
	n := strings.ToLower(name)
	for _, pat := range c.Exclude {
		if ok, _ := filepath.Match(strings.ToLower(pat), n); ok {
			return true
		}
	}
	return false
}

func (c *Config) normalize() {
	c.BaseDir = ExpandHome(strings.TrimSpace(c.BaseDir))
	if c.BaseDir != "" {
		if abs, err := filepath.Abs(c.BaseDir); err == nil {
			c.BaseDir = abs
		}
	}
	if c.Exclude == nil {
		c.Exclude = []string{}
	}
	clean := c.Exclude[:0]
	for _, e := range c.Exclude {
		if e = strings.TrimSpace(e); e != "" {
			clean = append(clean, e)
		}
	}
	c.Exclude = clean
}

// ExpandHome reemplaza un "~" inicial por la carpeta del usuario.
func ExpandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") || strings.HasPrefix(p, `~\`) {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[1:])
		}
	}
	return p
}
