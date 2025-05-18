package config

import (
	"dill-monitor/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Addresses []models.Address `json:"addresses"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(configPath string) (*Config, error) {
	// Ensure the config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		config := &Config{
			Addresses: []models.Address{},
		}
		if err := SaveConfig(configPath, config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		return config, nil
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to a JSON file
func SaveConfig(configPath string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// AddAddress adds a new address to the configuration
func (c *Config) AddAddress(address models.Address) error {
	// Check if address already exists
	for _, addr := range c.Addresses {
		if addr.Address == address.Address {
			return fmt.Errorf("address already exists")
		}
	}

	c.Addresses = append(c.Addresses, address)
	return nil
}

// RemoveAddress removes an address from the configuration
func (c *Config) RemoveAddress(address string) error {
	for i, addr := range c.Addresses {
		if addr.Address == address {
			c.Addresses = append(c.Addresses[:i], c.Addresses[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("address not found")
}

// GetAddress returns an address from the configuration
func (c *Config) GetAddress(address string) (*models.Address, error) {
	for _, addr := range c.Addresses {
		if addr.Address == address {
			return &addr, nil
		}
	}
	return nil, fmt.Errorf("address not found")
}

// ListAddresses returns all addresses in the configuration
func (c *Config) ListAddresses() []models.Address {
	return c.Addresses
}
