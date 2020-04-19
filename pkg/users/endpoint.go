package users

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type registerRequest struct {
	Input *model.UserInput
}

type registerResponse struct {
	User *model.UserOutput
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
	Input *model.UserInput
}

type loginResponse struct {
	User *model.UserOutput
	Err  error
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		user, err := s.Login(ctx, req.Input)
		return loginResponse{User: user, Err: err}, nil
	}
}

type refreshTokenRequest struct {
	Token string
}

type refreshTokenResponse struct {
	User *model.UserOutput
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
	Users []*model.User
	Err   error
}

func makeGetUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, err := s.GetUsers(ctx)
		return getUsersResponse{Users: users, Err: err}, nil
	}
}
