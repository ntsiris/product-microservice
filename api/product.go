package api

import (
	"net/http"
	"ntsiris/product-microservice/internal/utils"
)

type ProductHandler struct{}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// TODO: make assertion from the error to get the status code
			utils.WriteJSON(w, r.Response.StatusCode, APIError{Error: err.Error()})
		}
	}
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (handler *ProductHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /add", makeHTTPHandleFunc(handler.handleCreate))
}

func (handler *ProductHandler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler *ProductHandler) handleRetrieve(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler *ProductHandler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler *ProductHandler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	return nil
}
