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
	router.HandleFunc("POST /product/create", makeHTTPHandleFunc(handler.handleCreate))

	router.HandleFunc("GET /product/{id}", makeHTTPHandleFunc(handler.handleRetrieve))
	router.HandleFunc("GET /product", makeHTTPHandleFunc(handler.handleRetrieveAll))

	router.HandleFunc("PUT /product/update/", makeHTTPHandleFunc(handler.handleUpdate))

	router.HandleFunc("DELETE /product/delete/{id}", makeHTTPHandleFunc(handler.handleDelete))
}

func (handler *ProductHandler) handleCreate(w http.ResponseWriter, r *http.Request) error {
	productPayload := new(service.ProductCreationPayload)

	if err := parsePayload(r, productPayload); err != nil {
		return err
	}

	if err := validateStruct(r, productPayload); err != nil {
		return err
	}

	product := service.NewProduct(productPayload)
	err := handler.store.Create(&product)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "Product not created",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return utils.WriteJSON(w, http.StatusCreated, product)
}

func (handler *ProductHandler) handleRetrieve(w http.ResponseWriter, r *http.Request) error {
	requestedID, err := parseIntPathValue(r, "id")
	if err != nil {
		return err
	}

	requestedProduct, err := handler.retrieveProduct(r, service.ProductID(requestedID))
	if err != nil {
		return err
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
				Message:       "Invalid page number",
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
				Message:       "Invalid limit number",
				Operation:     types.FormatOperation(r.Method, r.URL.Path),
				EmbeddedError: err.Error(),
			}
		}
	}

	products, err := handler.store.RetrieveAll(page, limit)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "Error in product retrieval",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	if len(products) == 0 {
		products = []*service.Product{}
	}

	return utils.WriteJSON(w, http.StatusOK, products)
}

func (handler *ProductHandler) handleUpdate(w http.ResponseWriter, r *http.Request) error {
	updatePayload := service.NewDefaultUpdatePayload()
	if err := parsePayload(r, updatePayload); err != nil {
		return err
	}

	if err := validateStruct(r, updatePayload); err != nil {
		return err
	}

	product, err := handler.retrieveProduct(r, updatePayload.ID)
	if err != nil {
		return err
	}

	service.UpdateProduct(product, updatePayload)

	err = handler.store.Update(&product)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "Product not updated",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return utils.WriteJSON(w, http.StatusOK, product)
}

func (handler *ProductHandler) handleDelete(w http.ResponseWriter, r *http.Request) error {

	requestedID, err := parseIntPathValue(r, "id")
	if err != nil {
		return err
	}

	requestedProduct, err := handler.retrieveProduct(r, service.ProductID(requestedID))
	if err != nil {
		return err
	}

	err = handler.store.Delete(requestedProduct)
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "Product not deleted",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return utils.WriteJSON(w, http.StatusOK, requestedProduct)
}

func (handler *ProductHandler) retrieveProduct(r *http.Request, productID service.ProductID) (*service.Product, error) {
	requestedProduct, err := handler.store.Retrieve(service.ProductID(productID))
	if err != nil {
		return nil, &types.APIError{
			Code:          http.StatusNotFound,
			Message:       "Product not found",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return requestedProduct, nil
}

func parsePayload(r *http.Request, payload any) error {
	if err := utils.ParseJSON(r, payload); err != nil {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "Request Body parsing failed",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return nil
}

func validateStruct(r *http.Request, st any) error {
	if err := utils.Validate.Struct(st); err != nil {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "Request payload validation failed",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.(validator.ValidationErrors).Error(),
		}
	}

	return nil
}

func parseIntPathValue(r *http.Request, name string) (int64, error) {
	requestedValueStr := r.PathValue(name)

	value, err := strconv.ParseInt(requestedValueStr, 10, 64)
	if err != nil {
		return -1, &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "Invalid format of ID",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return value, nil
}
