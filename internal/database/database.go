package database

import (
	"database/sql"
	"fmt"
	"ntsiris/product-microservice/config"

	"github.com/go-sql-driver/mysql"
)

func NewSQLDatabase(config config.DBConfig) (*sql.DB, error) {
	dsn, err := getDSN(config)
	if err != nil {
		return nil, fmt.Errorf("create Data Source Name (DSN) from configuration: %v", err)
	}

	db, err := sql.Open(config.Driver, dsn)

	if err != nil {
		return nil, fmt.Errorf("initialize connection handle to MySQL database: %v", err)
	}

	return db, nil
}

func VerifyDatabaseConnection(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("establish connection to the Database: %v", err)
	}
	return nil
}

func getDSN(config config.DBConfig) (string, error) {
	switch config.Driver {
	case "mysql":
		mysqlConfig := mysql.Config{
			User:                 config.User,
			Passwd:               config.Password,
			Addr:                 config.Address,
			DBName:               config.Name,
			Net:                  "tcp",
			AllowNativePasswords: true,
			ParseTime:            true,
		}
		return mysqlConfig.FormatDSN(), nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
}
