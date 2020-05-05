package carts

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

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

func makeAddProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addProductRequest)
		setTmpUserIDCookie, userID, err := s.AddProduct(ctx, req.SKU, req.UserID, req.IsLoggedIn, req.IsTmpUserIDSet)
		return addProductResponse{Err: err, SetTmpUserIDCookie: setTmpUserIDCookie, UserID: userID}, nil
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
