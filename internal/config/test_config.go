package config

import (
	"path/filepath"
	"testing"
)

func setupTempHome(t *testing.T) string {
	t.Helper()
	tmp_dir := t.TempDir()
	t.Setenv("HOME", tmp_dir)
	return tmp_dir
}

func TestGetConfigPath(t *testing.T) {
	home := setupTempHome(t)

	got, err := getConfigPath()
	if err != nil {
		t.Fatalf("failed getting config path")
	}
	want := filepath.Join(home + "/" + configFileName)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
