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
	// The Product parameter might be modified
	Create(*Product) error

	// RetrieveAll returns a list of all available products.
	RetrieveAll(int, int) ([]*Product, error)

	// RetrieveOne returns a product by its unique ID.
	Retrieve(ProductID) (*Product, error)

	// UpdateOne updates the details of an existing product.
	// The Product parameter might be modified
	Update(*Product) error

	// DeleteOne removes a product from the store.
	Delete(*Product) error
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
