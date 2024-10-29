package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	t.Run("successfully creates a product with initial values", func(t *testing.T) {
		payload := &ProductCreationPayload{
			Price:       19.99,
			Quantity:    10,
			Discount:    5.0,
			Name:        "Test Product",
			Description: "A test product description",
		}

		product := NewProduct(payload)

		assert.Equal(t, payload.Price, product.Price)
		assert.Equal(t, payload.Quantity, product.Quantity)
		assert.Equal(t, float32(5.0), product.Discount) // Adjust type for Discount field
		assert.Equal(t, payload.Name, product.Name)
		assert.Equal(t, payload.Description, product.Description)
		assert.Equal(t, ProductID(0), product.ID)
		assert.WithinDuration(t, time.Now().UTC(), product.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), product.LastUpdated, time.Second)
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("updates product with new values", func(t *testing.T) {
		product := &Product{
			Price:       19.99,
			Quantity:    10,
			Discount:    5.0,
			Name:        "Original Product",
			Description: "Original description",
		}

		updatePayload := &ProductUpdatePayload{
			Price:       24.99,
			Quantity:    15,
			Discount:    10.0,
			Name:        "Updated Product",
			Description: "Updated description",
		}

		UpdateProduct(product, updatePayload)

		assert.Equal(t, updatePayload.Price, product.Price)
		assert.Equal(t, updatePayload.Quantity, product.Quantity)
		assert.Equal(t, float32(10.0), product.Discount)
		assert.Equal(t, updatePayload.Name, product.Name)
		assert.Equal(t, updatePayload.Description, product.Description)
		assert.Equal(t, 5, product.GetQuantityDelta())
		assert.WithinDuration(t, time.Now(), product.LastUpdated, time.Second)
	})

	t.Run("does not update fields with default values", func(t *testing.T) {
		product := &Product{
			Price:       19.99,
			Quantity:    10,
			Discount:    5.0,
			Name:        "Original Product",
			Description: "Original description",
		}

		updatePayload := &ProductUpdatePayload{
			Price:       -1,
			Quantity:    -1,
			Discount:    -1,
			Name:        "",
			Description: "",
		}

		UpdateProduct(product, updatePayload)

		assert.Equal(t, 19.99, product.Price)
		assert.Equal(t, 10, product.Quantity)
		assert.Equal(t, float32(5.0), product.Discount)
		assert.Equal(t, "Original Product", product.Name)
		assert.Equal(t, "Original description", product.Description)
	})

	t.Run("updates only specified fields", func(t *testing.T) {
		product := &Product{
			Price:       19.99,
			Quantity:    10,
			Discount:    5.0,
			Name:        "Original Product",
			Description: "Original description",
		}

		updatePayload := &ProductUpdatePayload{
			Name:        "Partially Updated Product",
			Description: "",
			Price:       25.99,
			Discount:    -1,
			Quantity:    -1,
		}

		UpdateProduct(product, updatePayload)

		assert.Equal(t, 25.99, product.Price)
		assert.Equal(t, "Partially Updated Product", product.Name)

		// These fields should not be updated
		assert.Equal(t, 10, product.Quantity)
		assert.Equal(t, float32(5.0), product.Discount)
		assert.Equal(t, "Original description", product.Description)
	})
}

func TestGetQuantityDelta(t *testing.T) {
	product := &Product{
		Quantity: 10,
	}

	updatePayload := &ProductUpdatePayload{
		Quantity: 15,
	}

	UpdateProduct(product, updatePayload)

	assert.Equal(t, 5, product.GetQuantityDelta())
}
