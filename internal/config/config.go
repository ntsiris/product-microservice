package config

// APIServerConfig holds the configuration settings for the API server.
type APIServerConfig struct {
	MigrateUp     bool
	MigrateDown   bool
	MigrationPath string
	PublicHost    string
	Port          string
	LogFile       string
	JWTSecret     string
}

// StorageConfig holds the configuration settings for the storage (database) connection.
type StorageConfig struct {
	Driver   string // Driver specifies the storage (database) driver type (e.g., "mysql" or "postgres").
	User     string // User is the username used for authenticating with the storage (database).
	Password string // Password is the password used for authenticating with the storage (database).
	Address  string // Address is the storage's (database) server address, including host and port.
	Name     string // Name is the name of the specific storage (database) to connect to.
}
