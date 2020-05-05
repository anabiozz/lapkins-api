package users

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

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := registerRequest{Input: &model.UserInput{}}
	if err := json.NewDecoder(r.Body).Decode(req.Input); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeRegisterResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(registerResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Path:    "/",
		Name:    "token",
		Value:   res.User.Token,
		Expires: res.User.ExpirationTime,
	})
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.User)
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := loginRequest{Input: &model.UserInput{}}
	if err := json.NewDecoder(r.Body).Decode(req.Input); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	tmpUserID, err := cookies.GetCookieValue(r, "tmp-user-id")
	if err != nil {
		if err != http.ErrNoCookie {
			return nil, err
		}
	}
	req.TmpUserID = tmpUserID
	return req, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(loginResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Path:    "/",
		Name:    "token",
		Value:   res.User.Token,
		Expires: res.User.ExpirationTime,
	})

	if res.UnsetTmpUserIDCookie {
		http.SetCookie(w, &http.Cookie{
			Path:    "/",
			Name:    "tmp-user-id",
			Value:   "",
			Expires: time.Unix(0, 0),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.User)
}

func decodeRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := refreshTokenRequest{}
	token, err := auth.GetToken(r)
	if err != nil {
		return nil, err
	}
	req.Token = token
	return req, nil
}

func encodeRefreshTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(refreshTokenResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.User)
}

func decodeGetUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	header := r.Header.Get("Authorization")
	_, err := auth.Check(header)
	if err != nil {
		return nil, errUnauthorized("%s", err)
	}
	return r, nil
}

func encodeGetUsersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getUsersResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Users)
}
