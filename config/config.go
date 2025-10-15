package config

import (
	"encoding/json"
	"fmt"
	"scp-go-identity-service/services"
	"os"
)

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig     `json:"database"`
	JWT      services.JWTConfig `json:"jwt"`
	Server   ServerConfig       `json:"server"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	ConnectionString string `json:"connectionString"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// LoadConfig loads configuration from a JSON file and overrides with environment variables.
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %v", err)
	}

	// Override sensitive data with environment variables if they are set.
	if jwtKey := os.Getenv("SCP_JWT_KEY"); jwtKey != "" {
		config.JWT.Key = jwtKey
	}

	if dbConnStr := os.Getenv("SCP_DB_CONNECTION_STRING"); dbConnStr != "" {
		config.Database.ConnectionString = dbConnStr
	}

	// Validate that critical configuration is not empty.
	if config.JWT.Key == "" {
		return nil, fmt.Errorf("JWT key is missing. Please set it in config.json or via the SCP_JWT_KEY environment variable")
	}
	if config.Database.ConnectionString == "" {
		return nil, fmt.Errorf("database connection string is missing. Please set it in config.json or via the SCP_DB_CONNECTION_STRING environment variable")
	}

	return &config, nil
}
