package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitAPIServerConfigFromEnv(t *testing.T) {
	// Reset environment variables to avoid interference from previous tests
	defer os.Unsetenv("PUBLIC_HOST")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("MIGRATE_UP")
	defer os.Unsetenv("MIGRATE_DOWN")
	defer os.Unsetenv("MIGRATION_PATH")
	defer os.Unsetenv("LOG_FILE")

	t.Run("environment variables are set", func(t *testing.T) {
		os.Setenv("PUBLIC_HOST", "testhost")
		os.Setenv("PORT", "9090")
		os.Setenv("MIGRATE_UP", "false")
		os.Setenv("MIGRATE_DOWN", "true")
		os.Setenv("MIGRATION_PATH", "test_migrations/")
		os.Setenv("LOG_FILE", "test.log")

		config := initAPIServerConfigFromEnv()

		assert.Equal(t, "testhost", config.PublicHost)
		assert.Equal(t, "9090", config.Port)
		assert.False(t, config.MigrateUp)
		assert.True(t, config.MigrateDown)
		assert.Equal(t, "test_migrations/", config.MigrationPath)
		assert.Equal(t, "test.log", config.LogFile)
	})

	t.Run("default values are applied when environment variables are missing", func(t *testing.T) {
		os.Unsetenv("PUBLIC_HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("MIGRATE_UP")
		os.Unsetenv("MIGRATE_DOWN")
		os.Unsetenv("MIGRATION_PATH")
		os.Unsetenv("LOG_FILE")

		config := initAPIServerConfigFromEnv()

		assert.Equal(t, "localhost", config.PublicHost)
		assert.Equal(t, "8080", config.Port)
		assert.True(t, config.MigrateUp)
		assert.False(t, config.MigrateDown)
		assert.Equal(t, "migrations/", config.MigrationPath)
		assert.Equal(t, "/var/log/product-api.log", config.LogFile)
	})
}

func TestInitDBConfigFromEnv(t *testing.T) {
	defer os.Unsetenv("DB_DRIVER")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("DB_NAME")

	t.Run("environment variables are set", func(t *testing.T) {
		os.Setenv("DB_DRIVER", "postgres")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_HOST", "testhost")
		os.Setenv("DB_PORT", "5432")
		os.Setenv("DB_NAME", "testDB")

		config := initStorageConfigFromEnv()

		assert.Equal(t, "postgres", config.Driver)
		assert.Equal(t, "testuser", config.User)
		assert.Equal(t, "testpass", config.Password)
		assert.Equal(t, "testhost:5432", config.Address)
		assert.Equal(t, "testDB", config.Name)
	})

	t.Run("default values are applied when environment variables are missing", func(t *testing.T) {
		os.Unsetenv("DB_DRIVER")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")

		config := initStorageConfigFromEnv()

		assert.Equal(t, "mysql", config.Driver)
		assert.Equal(t, "root", config.User)
		assert.Equal(t, "root", config.Password)
		assert.Equal(t, "localhost:3306", config.Address)
		assert.Equal(t, "productDB", config.Name)
	})
}

func TestGetEnvBool(t *testing.T) {
	t.Run("parses true values correctly", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "true")
		assert.True(t, getEnvBool("TEST_BOOL", false))

		os.Setenv("TEST_BOOL", "1")
		assert.True(t, getEnvBool("TEST_BOOL", false))

		os.Setenv("TEST_BOOL", "t")
		assert.True(t, getEnvBool("TEST_BOOL", false))

		os.Unsetenv("TEST_BOOL")
	})

	t.Run("parses false values correctly", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "false")
		assert.False(t, getEnvBool("TEST_BOOL", true))

		os.Setenv("TEST_BOOL", "0")
		assert.False(t, getEnvBool("TEST_BOOL", true))

		os.Setenv("TEST_BOOL", "f")
		assert.False(t, getEnvBool("TEST_BOOL", true))

		os.Unsetenv("TEST_BOOL")
	})

	t.Run("uses fallback when variable is not set", func(t *testing.T) {
		os.Unsetenv("TEST_BOOL")
		assert.True(t, getEnvBool("TEST_BOOL", true))
		assert.False(t, getEnvBool("TEST_BOOL", false))
	})
}
