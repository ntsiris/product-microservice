package storage

import (
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/service"
)

type MockProductStore struct {
	VerifyStoreFn    func() error
	VerifyStoreCalls int

	CreateFn    func(**service.Product) error
	CreateCalls int

	RetrieveAllFn    func(page, limit int) ([]*service.Product, error)
	RetrieveAllCalls int

	RetrieveFn    func(service.ProductID) (*service.Product, error)
	RetrieveCalls int

	UpdateFn    func(product **service.Product) error
	UpdateCalls int

	DeleteFn    func(product *service.Product) error
	DeleteCalls int
}

func (mock *MockProductStore) InitStore(*config.DBConfig) error {
	return nil
}

func (mock *MockProductStore) VerifyStoreConnection() error {
	mock.VerifyStoreCalls++
	return mock.VerifyStoreFn()
}

func (mock *MockProductStore) Create(product **service.Product) error {
	mock.CreateCalls++
	return mock.CreateFn(product)
}

func (mock *MockProductStore) RetrieveAll(page, limit int) ([]*service.Product, error) {
	mock.RetrieveAllCalls++
	return mock.RetrieveAllFn(page, limit)
}

func (mock *MockProductStore) Retrieve(id service.ProductID) (*service.Product, error) {
	mock.RetrieveCalls++
	return mock.RetrieveFn(id)
}

func (mock *MockProductStore) Update(product **service.Product) error {
	mock.UpdateCalls++
	return mock.UpdateFn(product)
}

func (mock *MockProductStore) Delete(product *service.Product) error {
	mock.DeleteCalls++
	return mock.DeleteFn(product)
}

func (mock *MockProductStore) RunMigrationUp(migrationPath string) error {
	return nil
}

func (mock *MockProductStore) RunMigrationDown(migrationPath string) error {
	return nil
}

func (mock *MockProductStore) Close() error {
	return nil
}
