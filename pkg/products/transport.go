package products

import (
	"context"
	"encoding/json"
	"net/http"
)

// encodeError writes a Service error to the given http.ResponseWriter.
func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	defaultErr := &Error{
		ErrorCode: http.StatusInternalServerError,
	}
	if err, ok := err.(*Error); ok {
		requestID, _ := getRequestID(ctx)
		defaultErr.ErrorCode = err.ErrorCode
		defaultErr.RequestID = requestID
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(defaultErr.ErrorCode)
	json.NewEncoder(w).Encode(defaultErr)
}

func decodeGetCatalogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getCatalogRequest{}
	req.Category = r.URL.Query().Get("category")
	return req, nil
}

func encodeGetCatalogResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getCatalogResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Products)
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

func decodeGetCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getCategoryRequest{}
	req.Category = r.URL.Query().Get("category")
	return req, nil
}

func encodeGetCategoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getCategoryResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Categories)
}

func decodeGetProductsByCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getProductsByCategoryRequest{}
	req.Category = r.URL.Query().Get("category")
	return req, nil
}

func encodeGetProductsByCategoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getProductsByCategoryResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Products)
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
