package config

// APIServerConfig holds the configuration settings for the API server.
type APIServerConfig struct {
	MigrateUp     bool   // MigrateUp indicates whether database migrations should run in the upward direction on startup.
	MigrateDown   bool   // MigrateDown indicates whether database migrations should run in the downward direction.
	MigrationPath string // MigrationPath specifies the path where migration files are located.
	PublicHost    string // PublicHost is the hostname or IP address where the API server is accessible.
	Port          string // Port is the network port on which the API server listens.
	LogFile       string // LogFile specifies the file path for storing server logs.
}

// StorageConfig holds the configuration settings for the storage (database) connection.
type StorageConfig struct {
	Driver   string // Driver specifies the storage (database) driver type (e.g., "mysql" or "postgres").
	User     string // User is the username used for authenticating with the storage (database).
	Password string // Password is the password used for authenticating with the storage (database).
	Address  string // Address is the storage's (database) server address, including host and port.
	Name     string // Name is the name of the specific storage (database) to connect to.
}
