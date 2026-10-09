package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed getting the home directory path: %w", err)
	}
	fullPath := filepath.Join(homeDir, configFileName)
	return fullPath, nil
}
func Read() (*Config, error) {
	resultConfig := &Config{}
	fullPath, err := getConfigPath()
	if err != nil {
		return &Config{}, err
	}
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to read gatorconfig.json: %w", err)
	}
	if err := json.Unmarshal(content, resultConfig); err != nil {
		return &Config{}, fmt.Errorf("failed unmarshaling json file: %w", err)
	}

	return resultConfig, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed unmarshaling config")
	}

	os.WriteFile(fullPath, data, 0644)
	return nil
}

func (cfg *Config) SetUser(newUserName string) error {
	cfg.CurrentUserName = newUserName
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}
