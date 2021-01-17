package erpsvc

import (
	"context"
	"github.com/anabiozz/core/lapkins/pkg/erp"
	"github.com/go-kit/kit/log"
)

type Storage interface {
	AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error)
	IncreaseProductQuantity(ctx context.Context, userID string, sku string) error
	DecreaseProductQuantity(ctx context.Context, userID string, sku string) error
	RemoveProduct(ctx context.Context, userID string, sku string) error
	LoadCart(ctx context.Context, userID string) ([]*erp.CartProduct, error)
	AddOrder(ctx context.Context, order *erp.Order) error
}

type Service interface {
	AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error)
	DecreaseProductQuantity(ctx context.Context, userID string, sku string) error
	IncreaseProductQuantity(ctx context.Context, userID string, sku string) error
	LoadCart(ctx context.Context, userID string) ([]*erp.CartProduct, error)
	RemoveProduct(ctx context.Context, userID string, sku string) error
}

type service struct {
	logger  log.Logger
	storage Storage
}

type ServiceConfig struct {
	Logger  log.Logger
	Storage Storage
}

func newService(cfg *ServiceConfig) (*service, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = log.NewNopLogger()
	}

	svc := &service{
		logger:  logger,
		storage: cfg.Storage,
	}

	return svc, nil
}

func (s *service) AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error) {
	setTmpUserIDCookie, userID, err := s.storage.AddProductToCard(ctx, sku, userID, isLoggedIn, isTmpUserIDSet)
	if err != nil {
		return false, "", erp.ErrBadRequest("%s", err)
	}
	return setTmpUserIDCookie, userID, nil
}

func (s *service) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	err := s.storage.DecreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		return erp.ErrBadRequest("%s", err)
	}
	return nil
}

func (s *service) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	err := s.storage.IncreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		return erp.ErrBadRequest("%s", err)
	}
	return nil
}

func (s *service) LoadCart(ctx context.Context, userID string) ([]*erp.CartProduct, error) {
	if userID == "" {
		return nil, erp.ErrBadRequest("%s", "provided user id is empty")
	}
	cart, err := s.storage.LoadCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return cart, nil

}

func (s *service) RemoveProduct(ctx context.Context, userID string, sku string) error {
	if userID == "" {
		return erp.ErrBadRequest("%s", "provided user id is empty")
	}
	err := s.storage.RemoveProduct(ctx, userID, sku)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddOrder(ctx context.Context, order *erp.Order) error {
	err := s.storage.AddOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}
