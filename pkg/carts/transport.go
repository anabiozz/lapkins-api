package carts

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/auth"
	"github.com/anabiozz/lapkins-api/pkg/cookies"
	"github.com/anabiozz/lapkins-api/pkg/model"
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

func decodeAddProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addProductRequest{User: &model.CartUser{}}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	token, err := cookies.GetCookieValue(r, "token")
	if err != nil {
		if err == http.ErrNoCookie {
			tmpCartSession, err := cookies.GetCookieValue(r, "tmp-cart-session")
			if err != nil {
				if err != http.ErrNoCookie {
					return nil, err
				}
			}
			req.User.TmpID = tmpCartSession
		} else {
			return nil, err
		}
	}
	if token != "" {
		claim, err := auth.Check(token)
		if err != nil {
			return nil, err
		}
		req.User.ID = claim.UserID
	}
	return req, nil
}

func encodeAddProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}

	if res.User.ID == "" && res.User.TmpID != "" {
		http.SetCookie(w, &http.Cookie{
			Path:    "/",
			Name:    "tmp-cart-session",
			Value:   res.User.TmpID,
			Expires: time.Now().Add(168 * time.Hour),
		})
	}

	if res.User.ID != "" && res.User.TmpID != "" {
		http.SetCookie(w, &http.Cookie{
			Path:    "/",
			Name:    "tmp-cart-session",
			Value:   "",
			Expires: time.Unix(0, 0),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

func decodeGetHeaderCartInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getHeaderCartInfoRequest{}
	userID, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeGetHeaderCartInfoResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getHeaderCartInfoResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Info)
}

func decodeGetCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getCartRequest{}
	userID, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
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
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	userID, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
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
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	userID, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
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

//
//func MakeHandler(ctx context.Context, s Service, st Storager) *mux.Router {
//	router := mux.NewRouter()
//	carts := router.PathPrefix("/carts").Subrouter()
//	carts.Handle("/add-product", Cors(s.AddProduct(ctx, st))).Methods("POST", "OPTIONS")
//	carts.Handle("/increase-product-quantity", Cors(s.IncreaseProductQuantity(st))).Methods("PUT", "OPTIONS")
//	carts.Handle("/decrease-product-quantity", Cors(s.DecreaseProductQuantity(st))).Methods("PUT", "OPTIONS")
//	carts.Handle("/remove-product", Cors(s.RemoveProduct(st))).Methods("DELETE", "OPTIONS")
//	carts.Handle("/load-carts", Cors(s.LoadCart(st))).Methods("GET", "OPTIONS")
//	carts.Handle("/create-order", Cors(s.CreateOrder(st))).Methods("POST", "OPTIONS")
//	return carts
//}
