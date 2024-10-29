package mocks

import (
	"errors"
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/service"
	"time"
)

// MockProductStore simulates the ProductStore interface for testing purposes.
type MockProductStore struct {
	Products map[int64]*service.Product // Simulates a database
	NextID   int64                      // Auto-increment ID for new products
	Err      error                      // Error to simulate failures
}

// NewMockProductStore initializes the mock with an empty product map.
func NewMockProductStore() *MockProductStore {
	return &MockProductStore{
		Products: make(map[int64]*service.Product),
		NextID:   1,
	}
}

// Create simulates adding a new product with auto-increment ID.
func (mock *MockProductStore) Create(product **service.Product) error {
	if mock.Err != nil {
		return mock.Err
	}
	(*product).ID = service.ProductID(mock.NextID)
	(*product).CreatedAt = time.Now()
	(*product).LastUpdated = time.Now()
	mock.Products[mock.NextID] = *product
	mock.NextID++
	return nil
}

// Retrieve finds a product by ID.
func (mock *MockProductStore) Retrieve(id service.ProductID) (*service.Product, error) {
	if mock.Err != nil {
		return nil, mock.Err
	}
	product, exists := mock.Products[int64(id)]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// RetrieveAll returns all products as a paginated result.
func (mock *MockProductStore) RetrieveAll(page, limit int) ([]*service.Product, error) {
	if mock.Err != nil {
		return nil, mock.Err
	}
	var products []*service.Product
	for _, product := range mock.Products {
		products = append(products, product)
	}
	return products, nil
}

// Update modifies an existing product's details.
func (mock *MockProductStore) Update(product **service.Product) error {
	if mock.Err != nil {
		return mock.Err
	}
	if _, exists := mock.Products[int64((*product).ID)]; !exists {
		return errors.New("product not found")
	}
	(*product).LastUpdated = time.Now() // Update the LastUpdated field
	mock.Products[int64((*product).ID)] = *product
	return nil
}

// Delete removes a product by ID.
func (mock *MockProductStore) Delete(product *service.Product) error {
	if mock.Err != nil {
		return mock.Err
	}
	if _, exists := mock.Products[int64(product.ID)]; !exists {
		return errors.New("product not found")
	}
	delete(mock.Products, int64(product.ID))
	return nil
}

// InitStore will not be tested since it can not be mocked.
func (mock *MockProductStore) InitStore(config *config.StorageConfig) error {
	return mock.Err
}

// VerifyStoreConnection will not be tested since it can not be mocked.
func (mock *MockProductStore) VerifyStoreConnection() error {
	return mock.Err
}

// Close will not be tested since it can not be mocked.
func (mock *MockProductStore) Close() error {
	return mock.Err
}

// RunMigrationUp will not be tested since it can not be mocked.
func (mock *MockProductStore) RunMigrationUp(path string) error {
	return mock.Err
}

// RunMigrationDown will not be tested since it can not be mocked.
func (mock *MockProductStore) RunMigrationDown(path string) error {
	return mock.Err
}
