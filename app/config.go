package main

import (
	"os"
)

// AppConfig represents the configuration for the app.
type AppConfig struct {
	VaultServiceName string
	VaultToken       string
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() (*AppConfig, error) {
	// Define the environment variables and their default values
	defaultEnvVars := map[string]string{
		"VAULT_SERVICE_NAME": "localhost",
		"VAULT_TOKEN":        "my-vault-token",
	}

	// Get the environment variable values or use the default values
	vaultServiceName := getEnv("VAULT_SERVICE_NAME", defaultEnvVars)
	vaultToken := getEnv("VAULT_TOKEN", defaultEnvVars)

	// Create the AppConfig instance
	config := &AppConfig{
		VaultServiceName: vaultServiceName,
		VaultToken:       vaultToken,
	}

	return config, nil
}

// getEnv retrieves the value of an environment variable.
// If the environment variable is not set, it returns the default value.
func getEnv(key string, defaultValues map[string]string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValues[key]
	}
	return value
}
