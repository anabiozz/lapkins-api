package carts

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
)

type Storage interface {
	AddProduct(ctx context.Context, sku string, user *model.CartUser) (*model.CartUser, error)
	IncreaseProductQuantity(ctx context.Context, userID string, sku string) error
	DecreaseProductQuantity(ctx context.Context, userID string, sku string) error
	RemoveProduct(variationID int, cartSession string, sizeOptionID int) error
	LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error)
	GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error)
}

type Service interface {
	AddProduct(ctx context.Context, sku string, user *model.CartUser) (*model.CartUser, error)
	CreateOrder(ctx context.Context) error
	DecreaseProductQuantity(ctx context.Context, userID string, sku string) error
	IncreaseProductQuantity(ctx context.Context, userID string, sku string) error
	LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error)
	RemoveProduct(ctx context.Context) error
	GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error)
}

type BasicService struct {
	logger  log.Logger
	storage Storage
}

type ServiceConfig struct {
	Logger  log.Logger
	Storage Storage
}

func NewService(cfg ServiceConfig) (*BasicService, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = log.NewNopLogger()
	}

	if cfg.Storage == nil {
		return nil, errBadRequest("storage must be provided")
	}

	svc := &BasicService{
		logger:  logger,
		storage: cfg.Storage,
	}

	return svc, nil
}

func (s *BasicService) GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error) {
	if userID == "" {
		return nil, errBadRequest("%s", "provided user id is empty")
	}
	info, err := s.storage.GetHeaderCartInfo(ctx, userID)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	return info, nil
}

func (s *BasicService) AddProduct(ctx context.Context, sku string, user *model.CartUser) (*model.CartUser, error) {
	user, err := s.storage.AddProduct(ctx, sku, user)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	return user, nil
}

func (s *BasicService) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *BasicService) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	err := s.storage.DecreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		return errBadRequest("%s", err)
	}
	return nil
}

func (s *BasicService) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	err := s.storage.IncreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		return errBadRequest("%s", err)
	}
	return nil
}

func (s *BasicService) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
	if userID == "" {
		return nil, errBadRequest("%s", "provided user id is empty")
	}
	cart, err := s.storage.LoadCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return cart, nil

}

func (s *BasicService) RemoveProduct(ctx context.Context) error {
	return nil

}
