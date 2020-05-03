package products

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
)

type Storage interface {
	GetCatalog(ctx context.Context, category string) ([]*model.CatalogProduct, error)
	GetProduct(ctx context.Context, sku string) (*model.Product, error)
	GetCategory(ctx context.Context, category string) ([]*model.Category, error)
	GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.Attribute) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string) error
	RemoveCategory(ctx context.Context, sku string) error
}

type Service interface {
	GetCategory(ctx context.Context, category string) ([]*model.Category, error)
	GetCatalog(ctx context.Context, category string) ([]*model.CatalogProduct, error)
	GetProduct(ctx context.Context, sku string) (*model.Product, error)
	GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.Attribute) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string) error
	RemoveCategory(ctx context.Context, sku string) error
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

func (s *BasicService) GetCategory(ctx context.Context, category string) ([]*model.Category, error) {
	categories, err := s.storage.GetCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *BasicService) GetCatalog(ctx context.Context, category string) ([]*model.CatalogProduct, error) {
	products, err := s.storage.GetCatalog(ctx, category)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *BasicService) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	product, err := s.storage.GetProduct(ctx, sku)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *BasicService) GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error) {
	products, err := s.storage.GetProductsByCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *BasicService) AddAttribute(ctx context.Context, sku string, attribute *model.Attribute) error {
	err := s.storage.AddAttribute(ctx, sku, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *BasicService) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	err := s.storage.RemoveAttribute(ctx, sku, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *BasicService) AddCategory(ctx context.Context, sku string) error {
	err := s.storage.AddCategory(ctx, sku)
	if err != nil {
		return err
	}
	return nil
}

func (s *BasicService) RemoveCategory(ctx context.Context, sku string) error {
	err := s.storage.RemoveCategory(ctx, sku)
	if err != nil {
		return err
	}
	return nil
}
