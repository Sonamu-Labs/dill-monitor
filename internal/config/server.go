package config

import (
	"dill-monitor/internal/models"
	"encoding/json"
	"os"
)

// LoadServerConfig loads server configuration from a file
func LoadServerConfig(path string) (*models.ServerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config models.ServerConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
