package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"ntsiris/product-microservice/internal/mocks"
	"ntsiris/product-microservice/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestProductHandler() (*ProductHandler, *mocks.MockProductStore) {
	mockStore := mocks.NewMockProductStore()
	handler := NewProductHandler(mockStore)
	return handler, mockStore
}

func TestHandleCreate(t *testing.T) {
	handler, mockStore := setupTestProductHandler()

	t.Run("successfully creates a product", func(t *testing.T) {
		payload := `{"name": "Test Product", "price": 100, "quantity": 10, "discount": 5.0, "description": "Test description"}`
		req := httptest.NewRequest(http.MethodPost, "/product/create", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleCreate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, int64(1), int64(mockStore.Products[1].ID))
		assert.Equal(t, "Test Product", mockStore.Products[1].Name)
		assert.Equal(t, "Test description", mockStore.Products[1].Description)
		assert.Equal(t, float64(100), float64(mockStore.Products[1].Price))
		assert.Equal(t, 10, mockStore.Products[1].Quantity)
		assert.Equal(t, float32(5.0), float32(mockStore.Products[1].Discount))
	})

	t.Run("returns 400 for invalid payload", func(t *testing.T) {
		payload := `{"name": ""}`
		req := httptest.NewRequest(http.MethodPost, "/product/create", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleCreate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 500 on store error", func(t *testing.T) {
		mockStore.Err = errors.New("db error")
		payload := `{"name": "Test Product", "price": 100, "quantity": 10, "discount": 5.0, "description": "Test description"}`
		req := httptest.NewRequest(http.MethodPost, "/product/create", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleCreate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockStore.Err = nil // Reset error for other tests
	})

	t.Run("fails with missing required fields", func(t *testing.T) {
		// Missing fields such as "price" and "quantity"
		payload := `{"name": "Test Product", "description": "This is a test"}`
		req := httptest.NewRequest(http.MethodPost, "/product/create", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleCreate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("fails with invalid data type for price", func(t *testing.T) {
		payload := `{"name": "Test Product", "price": "invalid_price", "quantity": 10, "description": "Test description"}`
		req := httptest.NewRequest(http.MethodPost, "/product/create", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleCreate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandleRetrieve(t *testing.T) {
	handler, mockStore := setupTestProductHandler()

	// Add a product to retrieve
	mockStore.Products[1] = &service.Product{ID: 1, Name: "Test Product"}

	t.Run("successfully retrieves a product", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product/1", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleRetrieve)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns 404 if product not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product/999", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleRetrieve)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestHandleRetrieveAll(t *testing.T) {
	handler, mockStore := setupTestProductHandler()

	// Populate mock store with test products
	mockStore.Products[1] = &service.Product{ID: 1, Name: "Product 1"}
	mockStore.Products[2] = &service.Product{ID: 2, Name: "Product 2"}

	t.Run("successfully retrieves all products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product?page=1&limit=10", nil)
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleRetrieveAll)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		// TODO: check the body here for the correct JSON response.
	})

	t.Run("returns empty list if no products found", func(t *testing.T) {
		// Clear products
		mockStore.Products = make(map[int64]*service.Product)

		req := httptest.NewRequest(http.MethodGet, "/product?page=1&limit=10", nil)
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleRetrieveAll)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[]`, rec.Body.String())
	})

	t.Run("returns 400 for invalid pagination parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/product?page=abc&limit=-1", nil)
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleRetrieveAll)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 500 on store error", func(t *testing.T) {
		mockStore.Err = errors.New("internal store error")
		req := httptest.NewRequest(http.MethodGet, "/product?page=1&limit=10", nil)
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleRetrieveAll)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockStore.Err = nil // Reset error for other tests
	})
}

func TestHandleUpdate(t *testing.T) {
	handler, mockStore := setupTestProductHandler()

	// Add a product to update
	mockStore.Products[1] = &service.Product{ID: 1, Name: "Original Product"}

	t.Run("successfully updates a product", func(t *testing.T) {
		payload := `{"id": 1, "name": "Updated Product", "price": 150, "quantity": 20, "description": "Updated description"}`
		req := httptest.NewRequest(http.MethodPut, "/product/update", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleUpdate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Updated Product", mockStore.Products[1].Name)
	})

	t.Run("returns 404 if product not found", func(t *testing.T) {
		payload := `{"id": 999, "name": "Non-existent Product"}`
		req := httptest.NewRequest(http.MethodPut, "/product/update", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleUpdate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("fails with non-existent product ID", func(t *testing.T) {
		payload := `{"id": 999, "name": "Updated Product"}`
		req := httptest.NewRequest(http.MethodPut, "/product/update", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleUpdate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("fails with invalid data type in update payload", func(t *testing.T) {
		payload := `{"id": 1, "price": "invalid_price", "name": "Updated Product"}`
		req := httptest.NewRequest(http.MethodPut, "/product/update", bytes.NewBufferString(payload))
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleUpdate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("partially updates product with valid fields", func(t *testing.T) {
		payload := `{"id": 1, "name": "Partially Updated Product"}`
		req := httptest.NewRequest(http.MethodPut, "/product/update", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleUpdate)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Partially Updated Product", mockStore.Products[1].Name)
	})
}

func TestHandleDelete(t *testing.T) {
	handler, mockStore := setupTestProductHandler()

	// Add a product to delete
	mockStore.Products[1] = &service.Product{ID: 1, Name: "Test Product"}

	t.Run("successfully deletes a product", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/product/delete/1", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleDelete)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		_, exists := mockStore.Products[1]
		assert.False(t, exists)
	})

	t.Run("returns 404 if product to delete not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/product/delete/999", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		handlerFunc := makeHTTPHandleFunc(handler.handleDelete)
		handlerFunc(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
