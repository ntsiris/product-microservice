package service

import (
	"time"
)

// ProductID is a unique identifier type for products.
type ProductID int64

// Product represents a product entity with details such as price, quantity, discount, and description.
type Product struct {
	Price         float64   `json:"price"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdated   time.Time `json:"lastUpdated"`
	ID            ProductID `json:"id"`
	Quantity      int       `json:"quantity"`
	quantityDelta int       // quantityDelta represents the change in quantity, used during updates.
	Discount      float32   `json:"discount"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
}

// ProductCreationPayload represents the required data to create a new product.
type ProductCreationPayload struct {
	Price       float64 `json:"price" validate:"required,number"`
	Quantity    int     `json:"quantity" validate:"required"`
	Discount    float32 `json:"discount"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
}

// ProductUpdatePayload represents the data used to update an existing product's details.
type ProductUpdatePayload struct {
	Price       float64   `json:"price" default:"-1" validate:"number"`
	ID          ProductID `json:"id" validate:"required"`
	Quantity    int       `json:"quantity"`
	Discount    float32   `json:"discount"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// ProductCRUDer defines an interface for CRUD operations
// on products, including create, retrieve, update, and delete methods.
type ProductCRUDer interface {
	// Create adds a new product to the store using the provided product reference.
	// The Product parameter may be modified with additional information (e.g., ID).
	Create(**Product) error

	// RetrieveAll retrieves a list of products with optional pagination.
	// The parameters specify the page and limit of products to retrieve.
	RetrieveAll(int, int) ([]*Product, error)

	// Retrieve fetches a product by its unique ID.
	Retrieve(ProductID) (*Product, error)

	// Update modifies the details of an existing product in the store.
	// The Product parameter may be modified with additional information.
	Update(**Product) error

	// Delete removes a specified product from the store.
	Delete(*Product) error
}

// NewProduct creates a new Product instance based on the provided ProductCreationPayload.
//
// Parameters:
// - productPayload: The payload containing product creation details.
//
// Returns:
// - A pointer to the newly created Product instance.
func NewProduct(productPayload *ProductCreationPayload) *Product {
	return &Product{
		Price:       productPayload.Price,
		CreatedAt:   time.Now().UTC(),
		LastUpdated: time.Now().UTC(),
		ID:          0,
		Quantity:    productPayload.Quantity,
		Discount:    productPayload.Discount,
		Name:        productPayload.Name,
		Description: productPayload.Description,
	}
}

// NewDefaultUpdatePayload creates a new ProductUpdatePayload with default values.
//
// Returns:
// - A pointer to a ProductUpdatePayload instance with default values indicating no update.
func NewDefaultUpdatePayload() *ProductUpdatePayload {
	return &ProductUpdatePayload{
		Price:       -1,
		ID:          0,
		Quantity:    -1,
		Discount:    -1,
		Name:        "",
		Description: "",
	}
}

// UpdateProduct applies updates to an existing Product based on the provided ProductUpdatePayload.
//
// Parameters:
// - product: The Product to update.
// - productUpdates: The ProductUpdatePayload containing update values.
//
// This function only updates fields in the Product that are set in the payload.
func UpdateProduct(product *Product, productUpdates *ProductUpdatePayload) {

	if productUpdates.Name != "" {
		product.Name = productUpdates.Name
	}

	if productUpdates.Description != "" {
		product.Description = productUpdates.Description
	}

	if productUpdates.Price >= 0 {
		product.Price = productUpdates.Price
	}

	if productUpdates.Quantity >= 0 {
		product.quantityDelta = productUpdates.Quantity - product.Quantity
		product.Quantity = productUpdates.Quantity
	}

	if productUpdates.Discount >= 0 {
		product.Discount = productUpdates.Discount
	}

	product.LastUpdated = time.Now()
}

// GetQuantityDelta returns the change in quantity during updates to a product.
func (product *Product) GetQuantityDelta() int {
	return product.quantityDelta
}
