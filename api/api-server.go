package api

import (
	"database/sql"
	"log"
	"net/http"
)

type APIError struct {
	Error string
}

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (server *APIServer) Run() error {
	router := http.NewServeMux()

	productHandler := NewProductHandler()
	productHandler.RegisterRoutes(router)

	subRouter := http.NewServeMux()
	subRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	log.Printf("Serving on address: %s\n", server.address)

	return http.ListenAndServe(server.address, subRouter)

}
