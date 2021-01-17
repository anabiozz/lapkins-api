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
		setTmpUserIDCookie, userID, err := s.AddProductToCard(ctx, req.SKU, req.UserID, req.IsLoggedIn, req.IsTmpUserIDSet)
		return addProductResponse{Err: err, SetTmpUserIDCookie: setTmpUserIDCookie, UserID: userID}, nil
	}
}

//type getHeaderCartInfoRequest struct {
//	UserID string `json:"user_id"`
//	Err    error  `json:"err"`
//}
//
//type getHeaderCartInfoResponse struct {
//	Info *erp.HeaderCartInfo `json:"info"`
//	Err  error               `json:"err"`
//}
//
//func makeGetHeaderCartInfo(s Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(getHeaderCartInfoRequest)
//		info, err := s.GetHeaderCartInfo(ctx, req.UserID)
//		return getHeaderCartInfoResponse{Err: err, Info: info}, nil
//	}
//}

type getCartRequest struct {
	UserID string `json:"user_id"`
	Err    error  `json:"err"`
}

type getCartResponse struct {
	Cart []*erp.CartProduct `json:"cart"`
	Err  error              `json:"err"`
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

//
//func decodeGetHeaderCartInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	req := getHeaderCartInfoRequest{}
//	userID, _, err := auth.GetUserID(r)
//	if err != nil {
//		return nil, erp.ErrBadRequest("%s", err)
//	}
//	req.UserID = userID
//	return req, nil
//}
//
//func encodeGetHeaderCartInfoResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
//	res := response.(getHeaderCartInfoResponse)
//	if res.Err != nil {
//		encodeError(ctx, res.Err, w)
//		return nil
//	}
//	w.Header().Set("Content-Type", "application/json")
//	return json.NewEncoder(w).Encode(res.Info)
//}

func decodeGetCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getCartRequest{}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, erp.ErrBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeGetCartResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getCartResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Cart)
}

func decodeIncreaseProductQtyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := increaseProductQtyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
	}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, erp.ErrBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeIncreaseProductQtyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(increaseProductQtyResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

func decodeDecreaseProductQtyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := decreaseProductQtyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
	}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, erp.ErrBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeDecreaseProductQtyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(decreaseProductQtyResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

func decodeRemoveProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := removeProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
	}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, erp.ErrBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeRemoveProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(removeProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

//
//type addOrderRequest struct {
//	Order *erp.Order
//}
//
//type addOrderResponse struct {
//	Err error `json:"err"`
//}
//
//func makeAddOrderEndpoint(s Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(addOrderRequest)
//		err = s.AddOrder(ctx, req.Order)
//		return addOrderResponse{Err: err}, nil
//	}
//}
//
//func decodeAddOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	req := addOrderRequest{}
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
//	}
//	return req, nil
//}
//
//func encodeAddOrderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
//	res := response.(addOrderResponse)
//	if res.Err != nil {
//		encodeError(ctx, res.Err, w)
//		return nil
//	}
//	w.Header().Set("Content-Type", "application/json")
//	return json.NewEncoder(w).Encode(true)
//}

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
