package orders

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
)

type Storage interface {
	AddOrder(ctx context.Context, order *model.Order) error
}

type Service interface {
	AddOrder(ctx context.Context, order *model.Order) error
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

func (s *BasicService) AddOrder(ctx context.Context, order *model.Order) error {
	err := s.storage.AddOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}
