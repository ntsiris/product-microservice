package service

import (
	"time"
)

type ProductID int64

// Product represents a product entity with details such as price, quantity, and description.
type Product struct {
	Price         float64   `json:"price"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdated   time.Time `json:"lastUpdated"`
	ID            ProductID `json:"id"`
	Quantity      int       `json:"quantity"`
	quantityDelta int
	Discount      float32 `json:"discount"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
}

// ProductCreationPayload represents the data required to create a new product.
type ProductCreationPayload struct {
	Price       float64 `json:"price" validate:"required,number"`
	Quantity    int     `json:"quantity" validate:"required"`
	Discount    float32 `json:"discount"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
}

type ProductUpdatePayload struct {
	Price       float64   `json:"price" default:"-1" validate:"number"`
	ID          ProductID `json:"id" validate:"required"`
	Quantity    int       `json:"quantity"`
	Discount    float32   `json:"discount"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// ProductCRUDer defines an interface for a products CRUD interface.
// It includes methods to create, retrieve, update, and delete products.
type ProductCRUDer interface {
	// Create creates a new product using the provided payload.
	// The Product parameter might be modified
	Create(**Product) error

	// RetrieveAll returns a list of all available products.
	RetrieveAll(int, int) ([]*Product, error)

	// Retrieve returns a product by its unique ID.
	Retrieve(ProductID) (*Product, error)

	// Update updates the details of an existing product.
	// The Product parameter might be modified
	Update(**Product) error

	// Delete removes a product from the store.
	Delete(*Product) error
}

func (product *Product) GetQuantityDelta() int {
	return product.quantityDelta
}

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
