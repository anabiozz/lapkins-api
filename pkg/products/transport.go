package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

// ***************** GET PRODUCT ********************

type getProductRequest struct {
	SKU string
}

type getProductResponse struct {
	Product *model.Product `json:"product"`
	Err     error          `json:"err"`
}

func decodeGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getProductRequest{}
	req.SKU = r.URL.Query().Get("sku")
	return req, nil
}

func encodeGetProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Product)
}

func makeGetProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProductRequest)
		product, err := s.GetProduct(ctx, req.SKU)
		fmt.Println(product)
		return getProductResponse{Err: err, Product: product}, nil
	}
}

// ***************** UPDATE PRODUCT ********************

type updateProductRequest struct {
	Product *model.Product
}

type updateProductResponse struct {
	Err error `json:"err"`
}

func decodeUpdateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := updateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeUpdateProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(updateProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func makeUpdateProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateProductRequest)
		err = s.UpdateProduct(ctx, req.Product)
		return updateProductResponse{Err: err}, nil
	}
}

// **************** GET PRODUCTS *******************

type getProductsRequest struct {
}

type getProductsResponse struct {
	Products []*model.Product `json:"products"`
	Err      error            `json:"err"`
}

func makeGetProductsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		products, err := s.GetProducts(ctx)
		return getProductsResponse{Err: err, Products: products}, nil
	}
}

func decodeGetProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getProductsRequest{}
	return req, nil
}

func encodeGetProductsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getProductsResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Products)
}

// **************************** GET CATEGORIES *************************

type getCategoriesRequest struct{}

type getCategoriesResponse struct {
	Categories []*model.Category `json:"categories"`
	Err        error             `json:"err"`
}

func makeGetCategoriesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		categories, err := s.GetCategories(ctx)
		return getCategoriesResponse{Err: err, Categories: categories}, nil
	}
}

func decodeGetCategoriesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getCategoriesRequest{}
	return req, nil
}

func encodeGetCategoriesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getCategoriesResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Categories)
}

// **************************** GET ATTRIBUTES *************************

type addAttributeRequest struct {
	SKU       string
	Attribute *model.NameValue
}

type addAttributeResponse struct {
	Err error `json:"err"`
}

func makeAddAttributeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addAttributeRequest)
		err = s.AddAttribute(ctx, req.SKU, req.Attribute)
		return addAttributeResponse{Err: err}, nil
	}
}

func decodeAddAttributeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addAttributeRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeAddAttributeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addAttributeResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

// **************************** REMOVE ATTRIBUTES *************************

type removeAttributeRequest struct {
	SKU       string
	Attribute string
}

type removeAttributeResponse struct {
	Err error `json:"err"`
}

func makeRemoveAttributeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeAttributeRequest)
		err = s.RemoveAttribute(ctx, req.SKU, req.Attribute)
		return removeAttributeResponse{Err: err}, nil
	}
}

func decodeRemoveAttributeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := removeAttributeRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeRemoveAttributeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(removeAttributeResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

// **************************** ADD CATEGORY *************************

type addCategoryRequest struct {
	SKU      string
	Category *model.Category
}

type addCategoryResponse struct {
	Err error `json:"err"`
}

func makeAddCategoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addCategoryRequest)
		err = s.AddCategory(ctx, req.SKU, req.Category)
		return addCategoryResponse{Err: err}, nil
	}
}

func decodeAddCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeAddCategoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addCategoryResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

// **************************** REMOVE CATEGORY *************************

type removeCategoryRequest struct {
	SKU      string
	Category *model.Category
}

type removeCategoryResponse struct {
	Err error `json:"err"`
}

func makeRemoveCategoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeCategoryRequest)
		err = s.RemoveCategory(ctx, req.SKU, req.Category)
		return removeCategoryResponse{Err: err}, nil
	}
}

func decodeRemoveCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := removeCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeRemoveCategoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(removeCategoryResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	defaultErr := &serviceError{code: http.StatusInternalServerError, message: "internal error"}
	if err, ok := err.(*serviceError); ok {
		defaultErr.code = err.code
		defaultErr.message = err.message
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(defaultErr.code)
	json.NewEncoder(w).Encode(defaultErr)
}
