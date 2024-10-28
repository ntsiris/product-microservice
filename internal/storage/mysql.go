package storage

import (
	"database/sql"
	"fmt"
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/service"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql" // MySQL driver
	_ "github.com/golang-migrate/migrate/v4/source/file"               // File source driver
)

type MySQLStore struct {
	db    *sql.DB
	dbURL string
}

const SQL_DRIVER string = "mysql"

func (mysqlStore *MySQLStore) InitStore(config *config.DBConfig) error {
	mysqlConfig := mysql.Config{
		User:   config.User,
		Passwd: config.Password,
		Addr:   config.Address,
		DBName: config.Name,
	}

	mysqlStore.dbURL = mysqlConfig.FormatDSN()
	db, err := sql.Open(SQL_DRIVER, mysqlStore.dbURL)

	if err != nil {
		return fmt.Errorf("failed to acquire storage connection handle: %v", err)
	}

	mysqlStore.db = db

	return nil
}

func (mysqlStore *MySQLStore) VerifyStoreConnection() error {
	err := mysqlStore.db.Ping()
	if err != nil {
		return fmt.Errorf("establish connection to the storage: %v", err)
	}

	return nil
}

func (mysqlStore *MySQLStore) Close() error {
	return mysqlStore.db.Close()
}

func (mysqlStore *MySQLStore) setUpMigration(migrationPath string) (*migrate.Migrate, error) {
	driver, err := mysqlMigrate.WithInstance(mysqlStore.db, &mysqlMigrate.Config{})
	if err != nil {
		return nil, fmt.Errorf("initialize migration driver: %v", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		SQL_DRIVER,
		driver)
	if err != nil {
		return nil, fmt.Errorf("establish connection to the database: %v", err)
	}

	return migrator, nil
}

func (mysqlStore *MySQLStore) RunMigrationUp(migrationPath string) error {
	migrator, err := mysqlStore.setUpMigration(migrationPath)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run up migrations: %v", err)
	}

	return nil
}

func (mysqlStore *MySQLStore) RunMigrationDown(migrationPath string) error {
	migrator, err := mysqlStore.setUpMigration(migrationPath)
	if err != nil {
		return err
	}

	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run down migrations: %v", err)
	}

	return nil
}

func (mysqlStore *MySQLStore) Create(product *service.Product) error {
	query := `INSERT INTO products (name, description, price, discount, quantity, createdAt, lastUpdated) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := mysqlStore.db.Exec(query, product.Name, product.Description, product.Price, product.Discount, product.Quantity, product.CreatedAt, product.LastUpdated)
	if err != nil {
		return err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("could not retrieve last inserted ID: %v", err)
	}

	product, err = mysqlStore.retrieveProductByID(service.ProductID(productID))
	if err != nil {
		return fmt.Errorf("failed to retrieve newly created product: %v", err)
	}

	return nil
}

func (mysqlStore *MySQLStore) RetrieveAll(page, limit int) ([]*service.Product, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	query := `SELECT * FROM products LIMIT ? OFFSET ?`

	rows, err := mysqlStore.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*service.Product
	var createdAtStr, lastUpdateStr string

	for rows.Next() {
		product := new(service.Product)
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Discount,
			&product.Quantity,
			&createdAtStr,
			&lastUpdateStr,
		)
		if err != nil {
			return nil, err
		}

		product.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to pase createdAt: %v", err)
		}

		product.LastUpdated, err = time.Parse("2006-01-02 15:04:05", lastUpdateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to pase lastUpdated: %v", err)
		}

		products = append(products, product)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (mysqlStore *MySQLStore) Retrieve(id service.ProductID) (*service.Product, error) {
	product, err := mysqlStore.retrieveProductByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (mysqlStore *MySQLStore) Update(product *service.Product) error {
	return nil
}

func (mysqlStore *MySQLStore) Delete(product *service.Product) error {
	return nil
}

func scanRowToProduct(row *sql.Row) (*service.Product, error) {
	product := new(service.Product)
	var createdAtStr, lastUpdateStr string

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Discount,
		&product.Quantity,
		&createdAtStr,
		&lastUpdateStr,
	)
	if err != nil {
		return nil, err
	}

	product.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to pase createdAt: %v", err)
	}

	product.LastUpdated, err = time.Parse("2006-01-02 15:04:05", lastUpdateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to pase lastUpdated: %v", err)
	}

	return product, nil
}

func (mysqlStore *MySQLStore) retrieveProductByID(id service.ProductID) (*service.Product, error) {
	query := `SELECT * FROM products WHERE id = ?`
	row := mysqlStore.db.QueryRow(query, id)

	product, err := scanRowToProduct(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return product, nil
}
