package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GITNOMAD_CONFIG", filepath.Join(dir, "sub", "config.json"))

	if _, exists, err := Load(); err != nil || exists {
		t.Fatalf("Load sin archivo: exists=%v err=%v", exists, err)
	}
	cfg := Default()
	cfg.BaseDir = filepath.Join(dir, "repos")
	cfg.Parallel = 8
	cfg.Exclude = []string{" test-* ", "", "Viejo"}
	if err := Save(cfg); err != nil {
		t.Fatal(err)
	}
	got, exists, err := Load()
	if err != nil || !exists {
		t.Fatalf("Load: exists=%v err=%v", exists, err)
	}
	if got.BaseDir != cfg.BaseDir || got.Parallel != 8 || len(got.Exclude) != 2 {
		t.Fatalf("config distinta: %+v", got)
	}
}

func TestValidate(t *testing.T) {
	base := Default()
	cases := map[string]func(*Config){
		"vacío":    func(c *Config) { c.BaseDir = "" },
		"relativo": func(c *Config) { c.BaseDir = "relativo/x" },
		"parallel": func(c *Config) { c.Parallel = 0 },
		"patrón":   func(c *Config) { c.Exclude = []string{"["} },
	}
	for name, mod := range cases {
		c := base
		mod(&c)
		if c.Validate() == nil {
			t.Errorf("%s: se esperaba error", name)
		}
	}
	if err := base.Validate(); err != nil {
		t.Errorf("default inválido: %v", err)
	}
}

func TestIsExcluded(t *testing.T) {
	c := Config{Exclude: []string{"test-*", "Viejo"}}
	for name, want := range map[string]bool{
		"test-api": true, "TEST-x": true, "viejo": true, "MejorarWindows": false, "mytest-a": false,
	} {
		if got := c.IsExcluded(name); got != want {
			t.Errorf("IsExcluded(%q)=%v, se esperaba %v", name, got, want)
		}
	}
}

func TestExpandHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	if got := ExpandHome("~/dev"); got != filepath.Join(home, "dev") {
		t.Errorf("ExpandHome: %q", got)
	}
	if got := ExpandHome("/abs"); got != "/abs" {
		t.Errorf("ExpandHome no debería tocar rutas absolutas: %q", got)
	}
}
