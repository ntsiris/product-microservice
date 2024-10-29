package storage

import (
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/service"
)

// ProductStore is an interface that extends the ProductCRUDer interface with additional methods
// for initializing, verifying, and managing the lifecycle of the product data store.
type ProductStore interface {
	service.ProductCRUDer // Embeds CRUD operations for managing product records.

	// InitStore initializes the connection to the product data store using the provided configuration.
	//
	// Parameters:
	// - config: A pointer to a StorageConfig struct that contains settings needed to initialize the store.
	//
	// Returns:
	// - An error if the store initialization fails; otherwise, nil.
	InitStore(*config.StorageConfig) error

	// VerifyStoreConnection checks the connection to the product data store to ensure it's accessible.
	//
	// Returns:
	// - An error if the store connection verification fails; otherwise, nil.
	VerifyStoreConnection() error

	// Close terminates the connection to the product data store, releasing any held resources.
	//
	// Returns:
	// - An error if the store closure fails; otherwise, nil.
	Close() error

	// RunMigrationUp applies database migrations in the upward direction, setting up or updating
	// the product database schema as necessary.
	//
	// Parameters:
	// - migrationPath: A string representing the file path to the migration scripts.
	//
	// Returns:
	// - An error if the migration fails; otherwise, nil.
	RunMigrationUp(string) error

	// RunMigrationDown rolls back database migrations, effectively reverting changes to the product
	// database schema.
	//
	// Parameters:
	// - migrationPath: A string representing the file path to the migration scripts.
	//
	// Returns:
	// - An error if the migration rollback fails; otherwise, nil.
	RunMigrationDown(string) error
}
