package api

import (
	"log"
	"net/http"
	"ntsiris/product-microservice/internal/storage"
)

// APIServer represents the server for handling API requests.
// It contains configuration for the server's address and a reference to the storage layer.
type APIServer struct {
	address string
	store   storage.ProductStore
}

// NewAPIServer initializes a new APIServer with the specified address and product storage layer.
//
// Parameters:
// - address: The network address the server listens on.
// - store: The storage layer used by the server to manage product data.
//
// Returns:
// - A pointer to the newly created APIServer instance.
func NewAPIServer(address string, store storage.ProductStore) *APIServer {
	return &APIServer{
		address: address,
		store:   store,
	}
}

// Run starts the API server, setting up routing and initializing the HTTP server.
//
// This method configures routing for the API, registers the product handler routes,
// and starts listening for incoming HTTP requests at the specified address.
//
// Returns:
// - An error if the server fails to start or encounters issues while running.
func (server *APIServer) Run() error {
	router := http.NewServeMux()

	productHandler := NewProductHandler(server.store)
	productHandler.RegisterRoutes(router)

	subRouter := http.NewServeMux()
	subRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	log.Printf("Product API Server running on address: %s\n", server.address)
	return http.ListenAndServe(server.address, subRouter)
}
