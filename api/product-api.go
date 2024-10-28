package api

import (
	"net/http"
	"ntsiris/product-microservice/internal/service"
	"ntsiris/product-microservice/internal/storage"
	"ntsiris/product-microservice/internal/types"
	"ntsiris/product-microservice/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	store storage.ProductStore
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			utils.WriteJSON(w, err.(*types.APIError).Code, err)
		}
	}
}

func NewProductHandler(userStore storage.ProductStore) *ProductHandler {
	return &ProductHandler{store: userStore}
}

func (handler *ProductHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /create", makeHTTPHandleFunc(handler.handleCreate))

	router.HandleFunc("GET /product/{id}", makeHTTPHandleFunc(handler.handleRetrieve))
	router.HandleFunc("GET /products", makeHTTPHandleFunc(handler.handleRetrieveAll))

	router.HandleFunc("PUT /update/{id}", makeHTTPHandleFunc(handler.handleUpdate))
	// router.HandleFunc("PUT /update-batch", makeHTTPHandleFunc(handler.handleUpdate))
}

func (handler *ProductHandler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	productPayload := new(service.ProductCreationPayload)

	if err := utils.ParseJSON(r, productPayload); err != nil {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "missing request body",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	if err := utils.Validate.Struct(productPayload); err != nil {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "request payload failed validation",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.(validator.ValidationErrors).Error(),
		}
	}

	product := service.NewProduct(productPayload)
	err := handler.store.Create(product)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "could not create product",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return utils.WriteJSON(w, http.StatusCreated, product)
}

func (handler *ProductHandler) handleRetrieve(w http.ResponseWriter, r *http.Request) error {
	requestedIDstr := r.PathValue("id")

	requestedID, err := strconv.ParseInt(requestedIDstr, 10, 64)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "invalid id format",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}
	requestedProduct, err := handler.store.Retrieve(service.ProductID(requestedID))
	if err != nil {
		return &types.APIError{
			Code:          http.StatusNotFound,
			Message:       "requested product was not found",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return utils.WriteJSON(w, http.StatusOK, requestedProduct)
}

func (handler *ProductHandler) handleRetrieveAll(w http.ResponseWriter, r *http.Request) error {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	var err error
	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			return &types.APIError{
				Code:          http.StatusBadRequest,
				Message:       "invalid page number",
				Operation:     types.FormatOperation(r.Method, r.URL.Path),
				EmbeddedError: err.Error(),
			}
		}
	}

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit < 1 {
			return &types.APIError{
				Code:          http.StatusBadRequest,
				Message:       "invalid limit number",
				Operation:     types.FormatOperation(r.Method, r.URL.Path),
				EmbeddedError: err.Error(),
			}
		}
	}

	products, err := handler.store.RetrieveAll(page, limit)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "failed to retrieve products",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	if products == nil {
		return &types.APIError{
			Code:      http.StatusNotFound,
			Message:   "no products found",
			Operation: types.FormatOperation(r.Method, r.URL.Path),
		}
	}

	return utils.WriteJSON(w, http.StatusOK, products)
}

func (handler *ProductHandler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler *ProductHandler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	return nil
}
