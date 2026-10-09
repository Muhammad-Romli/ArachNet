package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTempHome(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	return tmpDir
}

func TestGetConfigPath(t *testing.T) {
	home := setupTempHome(t)

	got, err := getConfigPath()
	if err != nil {
		t.Fatalf("failed getting config path: %v", err)
	}

	want := filepath.Join(home, configFileName)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRead(t *testing.T) {
	home := setupTempHome(t)

	jsonText := `{"db_url": "some/postgres/server", "current_user_name": "Justicar"}`
	err := os.WriteFile(filepath.Join(home, configFileName), []byte(jsonText), 0644)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	cfg, err := Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if cfg.DBURL != "some/postgres/server" {
		t.Errorf("DBURL = %q, want %q", cfg.DBURL, "some/postgres/server")
	}
	if cfg.CurrentUserName != "Justicar" {
		t.Errorf("CurrentUserName = %q, want %q", cfg.CurrentUserName, "Justicar")
	}
}

func TestWrite(t *testing.T) {
	home := setupTempHome(t)

	want := Config{DBURL: "some/postgres/server", CurrentUserName: "Justicar"}
	if err := write(want); err != nil {
		t.Fatalf("failed writing file: %v", err)
	}

	fullPath := filepath.Join(home, configFileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("file was not written: %v", err)
	}

	var got Config
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed unmarshaling json file: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestSetUser(t *testing.T) {
	setupTempHome(t)

	cfg := &Config{DBURL: "some/postgres/server"}
	if err := cfg.SetUser("Justicar"); err != nil {
		t.Fatalf("SetUser failed: %v", err)
	}

	// memory
	if cfg.CurrentUserName != "Justicar" {
		t.Errorf("in-memory CurrentUserName = %q, want %q", cfg.CurrentUserName, "Justicar")
	}

	// disk
	reloaded, err := Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if reloaded.CurrentUserName != "Justicar" {
		t.Errorf("persisted CurrentUserName = %q, want %q", reloaded.CurrentUserName, "Justicar")
	}
}
