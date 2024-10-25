package config

import (
	"fmt"
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type APIServerConfig struct {
	PublicHost string
	Port       string
}

type DBConfig struct {
	Driver   string
	User     string
	Password string
	Address  string
	Name     string
}

var EnvAPIServerConfig APIServerConfig
var EnvDBConfig DBConfig

func init() {
	log.Print("Initializing configuration from environment")
	godotenv.Load()
	EnvAPIServerConfig = initAPIServerConfig()
	EnvDBConfig = initDBConfig()
}

func initAPIServerConfig() APIServerConfig {
	return APIServerConfig{
		PublicHost: getEnv("PUBLIC_HOST", "localhost"),
		Port:       getEnv("PORT", "8080"),
	}
}

func initDBConfig() DBConfig {
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
