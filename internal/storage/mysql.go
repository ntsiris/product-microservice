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

// MySQLStore is a struct that provides methods for interacting with a MySQL database.
// It supports CRUD operations and schema migrations for product data.
type MySQLStore struct {
	db    *sql.DB
	dbURL string
}

// SQL_DRIVER is a constant that specifies the database driver used for MySQL.
const SQL_DRIVER string = "mysql"

// Create inserts a new product into the MySQL database and updates the provided product
// reference with the newly created product’s details.
//
// Parameters:
// - product: A double pointer to a Product instance, updated with additional data.
//
// Returns:
// - An error if the insertion fails; otherwise, nil.

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

// RetrieveAll retrieves a paginated list of products from the database.
//
// Parameters:
// - page: The page number for pagination (default is 1 if less than 1).
// - limit: The number of records per page (default is 10 if less than 1).
//
// Returns:
// - A slice of Product pointers and nil if successful.
// - An error if the retrieval fails.
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

// Retrieve fetches a product by its unique ID from the MySQL database.
//
// Parameters:
// - id: The unique ProductID of the product to retrieve.
//
// Returns:
// - A pointer to the retrieved Product and nil if successful.
// - An error if the product does not exist or retrieval fails.
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

// Update modifies an existing product’s details in the database.
//
// Parameters:
// - product: A double pointer to the Product instance containing the updated details.
//
// Returns:
// - An error if the update fails; otherwise, nil.
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

// Delete removes a product from the MySQL database.
//
// Parameters:
// - product: A pointer to the Product instance to delete.
//
// Returns:
// - An error if the deletion fails; otherwise, nil.
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

// InitStore initializes the MySQL store connection using the provided StorageConfig.
//
// Parameters:
// - config: A pointer to a StorageConfig struct containing database connection settings.
//
// Returns:
// - An error if the connection initialization fails; otherwise, nil.
func (mysqlStore *MySQLStore) InitStore(config *config.StorageConfig) error {
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

// VerifyStoreConnection verifies that the MySQL store connection is active and operational.
//
// Returns:
// - An error if the connection cannot be established; otherwise, nil.
func (mysqlStore *MySQLStore) VerifyStoreConnection() error {
	err := mysqlStore.db.Ping()
	if err != nil {
		return fmt.Errorf("error: could not establish connection to the storage: %v", err)
	}

	return nil
}

// Close terminates the MySQL store connection, releasing resources.
//
// Returns:
// - An error if the closure fails; otherwise, nil.
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

// RunMigrationUp applies upward database migrations to set up or update the database schema.
//
// Parameters:
// - migrationPath: A string path to the migration files.
//
// Returns:
// - An error if the migration fails; otherwise, nil.
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

// RunMigrationDown rolls back database migrations, effectively reverting changes to the database schema.
//
// Parameters:
// - migrationPath: A string path to the migration files.
//
// Returns:
// - An error if the migration rollback fails; otherwise, nil.
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

// scanIntoProduct scans the result rows into a Product instance.
//
// Parameters:
// - rows: A pointer to sql.Rows containing the product data.
//
// Returns:
// - A pointer to a populated Product instance and nil if successful.
// - An error if scanning fails.
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
