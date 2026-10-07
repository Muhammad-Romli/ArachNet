package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	resultConfig := &Config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, fmt.Errorf("failed getting the home directory path: %w", err)
	}
	fullPath := homeDir + "gatorconfig.json"
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to read gatorconfig.json: %w", err)
	}
	if err := json.Unmarshal(content, resultConfig); err != nil {
		return &Config{}, fmt.Errorf("failed unmarshaling json file: %w", err)
	}

	return resultConfig, nil
}
