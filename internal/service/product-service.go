package service

import "time"

type ProductID int64

// Product represents a product entity with details such as price, quantity, and description.
type Product struct {
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated"`
	ID          ProductID `json:"id"`
	Quantity    int       `json:"quantity"`
	Discount    float32   `json:"discount"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// ProductCreationPayload represents the data required to create a new product.
type ProductCreationPayload struct {
	Price       float64 `json:"price" validate:"required,number"`
	Quantity    int     `json:"quantity" validate:"required"`
	Discount    float32 `json:"discount"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
}

// ProductCRUDer defines an interface for a products CRUD interface.
// It includes methods to create, retrieve, update, and delete products.
type ProductCRUDer interface {
	// CreateOne creates a new product using the provided payload.
	CreateOne(*Product) (*Product, error)

	// CreateMany creates a batch of new products using the provided payloads.
	// On failure error is returned and a slice of Pointers to Products that
	// contains the products that were successfully created.
	CreateMany([]*Product) ([]*Product, error)

	// RetrieveAll returns a list of all available products.
	RetrieveAll() ([]*Product, error)

	// RetrieveOne returns a product by its unique ID.
	RetrieveOne(ProductID) (*Product, error)

	// RetrieveMany returns a list of products that match the given IDs.
	RetrieveMany([]ProductID) ([]*Product, error)

	// UpdateOne updates the details of an existing product.
	UpdateOne(*Product) error

	// UpdateMany updates the details of a batch of products existing in the store.
	UpdateMany([]*Product) error

	// DeleteOne removes a product from the store.
	DeleteOne(*Product) error

	// DeleteMany removes a batch of products from the store.
	DeleteMany([]*Product) error
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
