package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lpernett/godotenv"
)

// EnvAPIServerConfig holds the configuration settings for the API server, initialized from environment variables.
var EnvAPIServerConfig APIServerConfig

// EnvStorageConfig holds the database configuration settings, initialized from environment variables.
var EnvStorageConfig StorageConfig

// init loads environment variables, initializes configuration settings, and assigns values to EnvAPIServerConfig and EnvStorageConfig.
//
// This function is automatically invoked when the package is imported, loading environment variables from a .env file
// (if available) and setting the values in the configuration structs.
func init() {
	log.Print("Initializing configuration from environment")
	godotenv.Load()
	EnvAPIServerConfig = initAPIServerConfigFromEnv()
	EnvStorageConfig = initStorageConfigFromEnv()
}

// initAPIServerConfigFromEnv initializes the API server configuration using environment variables or fallback defaults.
//
// Returns:
// - An APIServerConfig struct with settings derived from environment variables or default values if not set.
func initAPIServerConfigFromEnv() APIServerConfig {
	return APIServerConfig{
		PublicHost:    getEnv("PUBLIC_HOST", "localhost"),
		Port:          getEnv("PORT", "8080"),
		MigrateUp:     getEnvBool("MIGRATE_UP", true),
		MigrateDown:   getEnvBool("MIGRATE_DOWN", false),
		MigrationPath: getEnv("MIGRATION_PATH", "migrations/"),
		LogFile:       getEnv("LOG_FILE", "/var/log/product-api.log"),
	}
}

// initStorageConfigFromEnv initializes the database configuration using environment variables or fallback defaults.
//
// Returns:
// - A StorageConfig struct with settings derived from environment variables or default values if not set.
func initStorageConfigFromEnv() StorageConfig {
	return StorageConfig{
		Driver:   getEnv("DB_DRIVER", "mysql"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Address:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3306")),
		Name:     getEnv("DB_NAME", "productDB"),
	}
}

// getEnv retrieves the value of an environment variable. If the variable is not set, it returns the provided fallback value.
//
// Parameters:
// - key: The name of the environment variable to retrieve.
// - fallback: The fallback value to return if the environment variable is not set.
//
// Returns:
// - A string containing the environment variable's value or the fallback if not set.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// getEnvBool retrieves a boolean value from an environment variable. If the variable is not set, it returns the provided fallback.
//
// Parameters:
// - key: The name of the environment variable to retrieve.
// - fallback: The fallback boolean value to return if the environment variable is not set.
//
// Returns:
// - A boolean indicating the environment variable's value or the fallback if not set.
func getEnvBool(key string, fallback bool) bool {
	if valStr, ok := os.LookupEnv(key); ok {
		valStr = strings.ToLower(valStr)
		return valStr == "true" || valStr == "t" || valStr == "1"
	}

	return fallback
}
