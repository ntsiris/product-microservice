package api

import (
	"database/sql"
	"net/http"
)

type APIError struct {
	Error string
}

type APIServer struct {
	Address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		Address: address,
		db:      db,
	}
}

func (server *APIServer) Run() error {
	router := http.NewServeMux()

	productHandler := NewProductHandler()
	productHandler.RegisterRoutes(router)

	subRouter := http.NewServeMux()
	subRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	return http.ListenAndServe(server.Address, subRouter)
}
