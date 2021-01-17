package erpsvc

import (
	"context"
	"encoding/json"
	"github.com/anabiozz/core/lapkins/pkg/auth"
	"github.com/anabiozz/core/lapkins/pkg/cookies"
	"github.com/anabiozz/core/lapkins/pkg/erp"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
)

// ***************** GET PRODUCT ********************

type getProductRequest struct {
	SKU string
}

type getProductResponse struct {
	Product *erp.Product `json:"product"`
	Err     error        `json:"err"`
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
		return getProductResponse{Err: err, Product: product}, nil
	}
}

// ***************** UPDATE PRODUCT ********************

type updateProductRequest struct {
	Product *erp.Product
}

type updateProductResponse struct {
	Err error `json:"err"`
}

func decodeUpdateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := updateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
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
	Products []*erp.Product `json:"products"`
	Err      error          `json:"err"`
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
	Categories []*erp.Category `json:"categories"`
	Err        error           `json:"err"`
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
	Attribute *erp.NameValue
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
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
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
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
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
	Category *erp.Category
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
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
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
	Category *erp.Category
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
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
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

type addProductRequest struct {
	SKU            string `json:"sku"`
	UserID         string `json:"user_id"`
	IsLoggedIn     bool   `json:"is_logged_in"`
	IsTmpUserIDSet bool   `json:"is_tmp_user_id_set"`
}

type addProductResponse struct {
	Err                error  `json:"err"`
	UserID             string `json:"user_id"`
	SetTmpUserIDCookie bool   `json:"set_tmp_user_id_cookie"`
}

func decodeAddProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
	}
	userID, isLoggedIn, err := auth.GetUserID(r)
	if err != nil {
		return nil, err
	}
	req.UserID = userID
	if userID != "" && isLoggedIn {
		req.IsLoggedIn = true
	}
	tmpUserID, err := cookies.GetCookieValue(r, "tmp-user-id")
	if err != nil {
		if err != http.ErrNoCookie {
			return nil, err
		}
	}
	if tmpUserID != "" {
		req.IsTmpUserIDSet = true
	}
	return req, nil
}

func encodeAddProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}

	if res.SetTmpUserIDCookie {
		http.SetCookie(w, &http.Cookie{
			Path:    "/",
			Name:    "tmp-user-id",
			Value:   res.UserID,
			Expires: time.Now().Add(168 * time.Hour),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

// ****************** Errors *********************

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	e, ok := err.(*erp.ServiceError)
	if !ok {
		e = &erp.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	e.Encode(w)
}

//
//func decodeError(r *http.Response) error {
//	e := &erp.ServiceError{}
//	e.Decode(r)
//
//	return e
//}
