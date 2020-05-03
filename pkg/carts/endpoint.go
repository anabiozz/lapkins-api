package carts

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type addProductRequest struct {
	SKU  string          `json:"sku"`
	User *model.CartUser `json:"user"`
}

type addProductResponse struct {
	User *model.CartUser `json:"user,omitempty"`
	Err  error           `json:"err"`
}

func makeAddProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addProductRequest)
		user, err := s.AddProduct(ctx, req.SKU, req.User)
		return addProductResponse{Err: err, User: user}, nil
	}
}

type getHeaderCartInfoRequest struct {
	UserID string `json:"user_id"`
	Err    error  `json:"err"`
}

type getHeaderCartInfoResponse struct {
	Info *model.HeaderCartInfo `json:"info"`
	Err  error                 `json:"err"`
}

func makeGetHeaderCartInfo(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getHeaderCartInfoRequest)
		info, err := s.GetHeaderCartInfo(ctx, req.UserID)
		return getHeaderCartInfoResponse{Err: err, Info: info}, nil
	}
}

type getCartRequest struct {
	UserID string `json:"user_id"`
	Err    error  `json:"err"`
}

type getCartResponse struct {
	Cart []*model.CartProduct `json:"cart"`
	Err  error                `json:"err"`
}

func makeGetCart(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCartRequest)
		cart, err := s.LoadCart(ctx, req.UserID)
		return getCartResponse{Err: err, Cart: cart}, nil
	}
}

type increaseProductQtyRequest struct {
	UserID string `json:"user_id"`
	SKU    string `json:"sku"`
}

type increaseProductQtyResponse struct {
	Err error `json:"err"`
}

func makeIncreaseProductQty(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(increaseProductQtyRequest)
		err = s.IncreaseProductQuantity(ctx, req.UserID, req.SKU)
		return increaseProductQtyResponse{Err: err}, nil
	}
}

type decreaseProductQtyRequest struct {
	UserID string `json:"user_id"`
	SKU    string `json:"sku"`
}

type decreaseProductQtyResponse struct {
	Err error `json:"err"`
}

func makeDecreaseProductQty(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(decreaseProductQtyRequest)
		err = s.DecreaseProductQuantity(ctx, req.UserID, req.SKU)
		return decreaseProductQtyResponse{Err: err}, nil
	}
}

type removeProductRequest struct {
	UserID string `json:"user_id"`
	SKU    string `json:"sku"`
}

type removeProductResponse struct {
	Err error `json:"err"`
}

func makeRemoveProduct(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeProductRequest)
		err = s.RemoveProduct(ctx, req.UserID, req.SKU)
		return removeProductResponse{Err: err}, nil
	}
}
