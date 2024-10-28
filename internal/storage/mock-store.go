package storage

import (
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/service"
)

type MockProductStore struct {
	VerifyStoreFn    func() error
	VerifyStoreCalls int

	CreateOneFn    func(*service.Product) (*service.Product, error)
	CreateOneCalls int

	CreateManeFn    func([]*service.Product) ([]*service.Product, error)
	CreateManyCalls int

	RetrieveAllFn    func() ([]*service.Product, error)
	RetrieveAllCalls int

	RetrieveOneFn    func(service.ProductID) (*service.Product, error)
	RetrieveOneCalls int

	RetrieveManyFn    func([]service.ProductID) ([]*service.Product, error)
	RetrieveManyCalls int

	UpdateOneFn    func(product *service.Product) error
	UpdateOneCalls int

	UpdateManyFn    func(products []*service.Product) error
	UpdateManyCalls int

	DeleteOneFn    func(product *service.Product) error
	DeleteOneCalls int

	DeleteManyFn    func(products []*service.Product) error
	DeleteManyCalls int
}

func (mock *MockProductStore) InitStore(*config.DBConfig) error {
	return nil
}

func (mock *MockProductStore) VerifyStoreConnection() error {
	mock.VerifyStoreCalls++
	return mock.VerifyStoreFn()
}

func (mock *MockProductStore) CreateOne(product *service.Product) (*service.Product, error) {
	mock.CreateOneCalls++
	return mock.CreateOneFn(product)
}

func (mock *MockProductStore) CreateMany(products []*service.Product) ([]*service.Product, error) {
	mock.CreateManyCalls++
	return mock.CreateManeFn(products)
}

func (mock *MockProductStore) RetrieveAll() ([]*service.Product, error) {
	mock.RetrieveAllCalls++
	return mock.RetrieveAllFn()
}

func (mock *MockProductStore) RetrieveOne(id service.ProductID) (*service.Product, error) {
	mock.RetrieveOneCalls++
	return mock.RetrieveOneFn(id)
}

func (mock *MockProductStore) RetrieveMany(id []service.ProductID) ([]*service.Product, error) {
	mock.RetrieveManyCalls++
	return mock.RetrieveManyFn(id)
}

func (mock *MockProductStore) UpdateOne(product *service.Product) error {
	mock.UpdateOneCalls++
	return mock.UpdateOneFn(product)
}

func (mock *MockProductStore) UpdateMany(products []*service.Product) error {
	mock.UpdateManyCalls++
	return mock.UpdateManyFn(products)
}
func (mock *MockProductStore) DeleteOne(product *service.Product) error {
	mock.DeleteOneCalls++
	return mock.DeleteOneFn(product)
}

func (mock *MockProductStore) DeleteMany(products []*service.Product) error {
	mock.DeleteManyCalls++
	return mock.DeleteManyFn(products)
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
