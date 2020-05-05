package orders

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type addOrderRequest struct {
	Order *model.Order
}

type addOrderResponse struct {
	Err error `json:"err"`
}

func makeAddOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addOrderRequest)
		err = s.AddOrder(ctx, req.Order)
		return addOrderResponse{Err: err}, nil
	}
}
