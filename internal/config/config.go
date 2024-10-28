package config

type APIServerConfig struct {
	MigrateUp     bool
	MigrateDown   bool
	MigrationPath string
	PublicHost    string
	Port          string
}

type DBConfig struct {
	Driver   string
	User     string
	Password string
	Address  string
	Name     string
}
