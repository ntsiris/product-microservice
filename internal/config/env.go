package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lpernett/godotenv"
)

var EnvAPIServerConfig APIServerConfig
var EnvDBConfig DBConfig

func init() {
	log.Print("Initializing configuration from environment")
	godotenv.Load()
	EnvAPIServerConfig = initAPIServerConfigFromEnv()
	EnvDBConfig = initDBConfigFromEnv()
}

func initAPIServerConfigFromEnv() APIServerConfig {
	return APIServerConfig{
		PublicHost:    getEnv("PUBLIC_HOST", "localhost"),
		Port:          getEnv("PORT", "8080"),
		MigrateUp:     getEnvBool("MIGRATE_UP", true),
		MigrateDown:   getEnvBool("MIGRATE_DOWN", false),
		MigrationPath: getEnv("MIGRATION_PATH", "migrations/"),
	}
}

func initDBConfigFromEnv() DBConfig {
	return DBConfig{
		Driver:   getEnv("DB_DRIVER", "mysql"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Address:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "localhost"), getEnv("DB_PORT", "3306")),
		Name:     getEnv("DB_NAME", "productDB"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if valStr, ok := os.LookupEnv(key); ok {
		valStr = strings.ToLower(valStr)
		return valStr == "true" || valStr == "t" || valStr == "1"
	}

	return fallback
}
