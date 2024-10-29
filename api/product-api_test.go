package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"ntsiris/product-microservice/internal/service"
	"ntsiris/product-microservice/internal/storage"
	"testing"
)

// Test handleCreate with successful product creation
func TestHandleCreate_Success(t *testing.T) {
	mockStore := &storage.MockProductStore{
		CreateFn: func(p **service.Product) error {
			return nil
		},
	}

	handler := NewProductHandler(mockStore)

	productPayload := &service.ProductCreationPayload{
		Name:        "A Test Product",
		Description: "A Product used for testing",
		Price:       0.99,
		Quantity:    1,
	}

	payloadBytes, _ := json.Marshal(productPayload)
	request := httptest.NewRequest("POST", "/add", bytes.NewBuffer(payloadBytes))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	makeHTTPHandleFunc(handler.handleCreate)(recorder, request)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("expected status code 201, got %d", status)
	}

	if mockStore.CreateCalls != 1 {
		t.Errorf("expected CreatedOne to be called once, got %d", mockStore.CreateCalls)
	}
}

// Test handleCreate with missing request body (BadRequest)
func TestHandleCreate_BadRequest(t *testing.T) {
	mockStore := &storage.MockProductStore{
		CreateFn: func(p **service.Product) error {
			return nil // Should not be called in this test
		},
	}

	handler := NewProductHandler(mockStore)

	// Prepare payload with empty body:
	request := httptest.NewRequest("POST", "/add", nil)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	makeHTTPHandleFunc(handler.handleCreate)(recorder, request)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", status)
	}

	if mockStore.CreateCalls != 0 {
		t.Errorf("expected CreateOne not to be called, got %d", mockStore.CreateCalls)
	}
}

// Test handleRetrieveOne with a successful product retrieval
func TestRetrieveOne(t *testing.T) {
	mockStore := &storage.MockProductStore{
		RetrieveFn: func(p service.ProductID) (*service.Product, error) {
			return &service.Product{
				ID:          1,
				Name:        "Test Product",
				Description: "A product for testing",
				Price:       0.99,
			}, nil
		},
	}

	handler := NewProductHandler(mockStore)

	request := httptest.NewRequest("GET", "/get/", nil)
	request.SetPathValue("id", "1")
	recorder := httptest.NewRecorder()

	makeHTTPHandleFunc(handler.handleRetrieve)(recorder, request)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code 200, got %d", status)
	}

	if mockStore.RetrieveCalls != 1 {
		t.Errorf("expected RetrieveOne to be called once, got %d", mockStore.RetrieveCalls)
	}
}

func TestHandleRetrieveOne_NotFound(t *testing.T) {
	mockStore := &storage.MockProductStore{
		RetrieveFn: func(p service.ProductID) (*service.Product, error) {
			return nil, fmt.Errorf("no product found") // Simulate no product found
		},
	}

	handler := NewProductHandler(mockStore)

	request := httptest.NewRequest("GET", "/get/", nil)
	request.SetPathValue("id", "99")
	recorder := httptest.NewRecorder()

	makeHTTPHandleFunc(handler.handleRetrieve)(recorder, request)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", status)
	}

	if mockStore.RetrieveCalls != 1 {
		t.Errorf("expected RetrieveOne to be called once, got %d", mockStore.CreateCalls)
	}
}
