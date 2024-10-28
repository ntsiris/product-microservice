package storage

import (
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/service"
)

type ProductStore interface {
	service.ProductCRUDer
	InitStore(*config.DBConfig) error
	VerifyStoreConnection() error
	RunMigrationUp(string) error
	RunMigrationDown(string) error
	Close() error
}
