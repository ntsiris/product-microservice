package api

import (
	"log"
	"net/http"
	"ntsiris/product-microservice/internal/storage"
)

type APIServer struct {
	address string
	store   storage.ProductStore
}

func NewAPIServer(address string, store storage.ProductStore) *APIServer {
	return &APIServer{
		address: address,
		store:   store,
	}
}

func (server *APIServer) Run() error {
	router := http.NewServeMux()

	productHandler := NewProductHandler(server.store)
	productHandler.RegisterRoutes(router)

	subRouter := http.NewServeMux()
	subRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	log.Printf("Product API Server running on address: %s\n", server.address)
	return http.ListenAndServe(server.address, subRouter)
}
