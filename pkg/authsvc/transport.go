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

type registerRequest struct {
	Input *erp.UserInput
}

type registerResponse struct {
	User *erp.UserOutput
	Err  error
}

func makeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		user, err := s.Register(ctx, req.Input)
		return registerResponse{User: user, Err: err}, nil
	}
}

type loginRequest struct {
	Input     *erp.UserInput
	TmpUserID string `json:"tmp_user_id"`
}

type loginResponse struct {
	User                 *erp.UserOutput
	UnsetTmpUserIDCookie bool `json:"unset_tmp_user_id_cookie"`
	Err                  error
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		user, unsetTmpUserIDCookie, err := s.Login(ctx, req.Input, req.TmpUserID)
		return loginResponse{User: user, UnsetTmpUserIDCookie: unsetTmpUserIDCookie, Err: err}, nil
	}
}

type refreshTokenRequest struct {
	Token string
}

type refreshTokenResponse struct {
	User *erp.UserOutput
	Err  error
}

func makeRefreshTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(refreshTokenRequest)
		user, err := s.RefreshToken(ctx, req.Token)
		return refreshTokenResponse{User: user, Err: err}, nil
	}
}

type getUsersResponse struct {
	Users []*erp.User
	Err   error
}

func makeGetUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, err := s.GetUsers(ctx)
		return getUsersResponse{Users: users, Err: err}, nil
	}
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := registerRequest{Input: &erp.UserInput{}}
	if err := json.NewDecoder(r.Body).Decode(req.Input); err != nil {
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
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
	req := loginRequest{Input: &erp.UserInput{}}
	if err := json.NewDecoder(r.Body).Decode(req.Input); err != nil {
		return nil, erp.ErrBadRequest("failed to decode JSON request: %v", err)
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
		return nil, erp.ErrUnauthorized("%s", err)
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
