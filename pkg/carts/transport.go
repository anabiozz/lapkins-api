package carts

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/auth"
	"github.com/anabiozz/lapkins-api/pkg/cookies"
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
	req := addProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
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

func decodeGetHeaderCartInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getHeaderCartInfoRequest{}
	userID, _, err := auth.GetUserID(r)
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
	userID, _, err := auth.GetUserID(r)
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
	userID, _, err := auth.GetUserID(r)
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
	userID, _, err := auth.GetUserID(r)
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

func decodeRemoveProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := removeProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
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
