package products

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type getCatalogRequest struct {
	Department string
	Category   string
}

type getCatalogResponse struct {
	Products []*model.CatalogProduct `json:"products"`
	Err      error                   `json:"err"`
}

func makeGetCatalogEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCatalogRequest)
		products, err := s.GetCatalog(ctx, req.Department, req.Category)
		return getCatalogResponse{Err: err, Products: products}, nil
	}
}

type getProductRequest struct {
	SKU  string
	Attr string
}

type getProductResponse struct {
	Product *model.VariationProduct `json:"product"`
	Err     error                   `json:"err"`
}

func makeGetProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProductRequest)
		product, err := s.GetProduct(ctx, req.SKU, req.Attr)
		return getProductResponse{Err: err, Product: product}, nil
	}
}

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

type getProductsByCategoryRequest struct {
	Category string
}

type getProductsByCategoryResponse struct {
	Products []*model.SKUProduct `json:"products"`
	Err      error               `json:"err"`
}

func makeGetProductsByCategoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProductsByCategoryRequest)
		products, err := s.GetProductsByCategory(ctx, req.Category)
		return getProductsByCategoryResponse{Err: err, Products: products}, nil
	}
}

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
