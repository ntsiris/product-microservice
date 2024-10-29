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

func (mysqlStore *MySQLStore) Create(product **service.Product) error {
	query := `INSERT INTO products (name, description, price, discount, quantity, createdAt, lastUpdated) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := mysqlStore.db.Exec(query,
		(*product).Name,
		(*product).Description,
		(*product).Price,
		(*product).Discount,
		(*product).Quantity,
		(*product).CreatedAt,
		(*product).LastUpdated)
	if err != nil {
		return err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error: could not retrieve last inserted ID: %v", err)
	}

	*product, err = mysqlStore.Retrieve(service.ProductID(productID))
	if err != nil {
		return fmt.Errorf("error: could not retrieve newly created product: %v", err)
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
	for rows.Next() {
		product, err := scanIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (mysqlStore *MySQLStore) Retrieve(id service.ProductID) (*service.Product, error) {
	query := `SELECT * FROM products WHERE id = ?`
	rows, err := mysqlStore.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoProduct(rows)
	}

	return nil, fmt.Errorf("error: product with id %d not found", id)
}

func (mysqlStore *MySQLStore) Update(product **service.Product) error {
	// Atomic increment of quantity field
	query := `UPDATE products SET name = ?, description = ?, price = ?, discount = ?, quantity = quantity + ?, lastUpdated = ? WHERE id = ? AND quantity + ? >= 0`

	result, err := mysqlStore.db.Exec(query,
		(*product).Name,
		(*product).Description,
		(*product).Price,
		(*product).Discount,
		(*product).GetQuantityDelta(),
		(*product).LastUpdated,
		(*product).ID,
		(*product).GetQuantityDelta())

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 1 {
		return fmt.Errorf("error: more than one rows were affected. rows affected: %d", rowsAffected)
	}

	*product, err = mysqlStore.Retrieve((*product).ID)
	if err != nil {
		return fmt.Errorf("error: could not retrieve updated product: %v", err)

	}

	return nil
}

func (mysqlStore *MySQLStore) Delete(product *service.Product) error {
	query := `DELETE FROM products WHERE id = ?`

	result, err := mysqlStore.db.Exec(query, product.ID)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (mysqlStore *MySQLStore) InitStore(config *config.DBConfig) error {
	mysqlConfig := mysql.Config{
		User:                 config.User,
		Passwd:               config.Password,
		Addr:                 config.Address,
		DBName:               config.Name,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.UTC,
	}

	mysqlStore.dbURL = mysqlConfig.FormatDSN()
	db, err := sql.Open(SQL_DRIVER, mysqlStore.dbURL)

	if err != nil {
		return fmt.Errorf("error: could not acquire storage connection handle: %v", err)
	}

	mysqlStore.db = db

	return nil
}

func (mysqlStore *MySQLStore) VerifyStoreConnection() error {
	err := mysqlStore.db.Ping()
	if err != nil {
		return fmt.Errorf("error: could not establish connection to the storage: %v", err)
	}

	return nil
}

func (mysqlStore *MySQLStore) Close() error {
	return mysqlStore.db.Close()
}

func (mysqlStore *MySQLStore) setUpMigration(migrationPath string) (*migrate.Migrate, error) {
	driver, err := mysqlMigrate.WithInstance(mysqlStore.db, &mysqlMigrate.Config{})
	if err != nil {
		return nil, fmt.Errorf("error: could not initialize migration driver: %v", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		SQL_DRIVER,
		driver)
	if err != nil {
		return nil, fmt.Errorf("error: could not establish connection to the database: %v", err)
	}

	return migrator, nil
}

func (mysqlStore *MySQLStore) RunMigrationUp(migrationPath string) error {
	migrator, err := mysqlStore.setUpMigration(migrationPath)
	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error: could not run up migrations: %v", err)
	}

	return nil
}

func (mysqlStore *MySQLStore) RunMigrationDown(migrationPath string) error {
	migrator, err := mysqlStore.setUpMigration(migrationPath)
	if err != nil {
		return err
	}

	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error: run down migrations failed: %v", err)
	}

	return nil
}

func scanIntoProduct(rows *sql.Rows) (*service.Product, error) {
	product := new(service.Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Discount,
		&product.Quantity,
		&product.CreatedAt,
		&product.LastUpdated,
	)

	return product, err
}
