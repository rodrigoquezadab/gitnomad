package gitops

import "testing"

func TestNormalizeURL(t *testing.T) {
	want := "github.com/rodrigoquezadab/mejorarwindows"
	for _, u := range []string{
		"https://github.com/rodrigoquezadab/MejorarWindows",
		"https://github.com/rodrigoquezadab/MejorarWindows.git",
		"https://github.com/rodrigoquezadab/MejorarWindows/",
		"https://token@github.com/rodrigoquezadab/MejorarWindows.git",
		"git@github.com:rodrigoquezadab/MejorarWindows.git",
		"ssh://git@github.com/rodrigoquezadab/MejorarWindows.git",
	} {
		if got := NormalizeURL(u); got != want {
			t.Errorf("NormalizeURL(%q) = %q, se esperaba %q", u, got, want)
		}
	}
	if NormalizeURL("https://github.com/a/b") == NormalizeURL("https://github.com/a/c") {
		t.Error("repos distintos no deberían coincidir")
	}
}
