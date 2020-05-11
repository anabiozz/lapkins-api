package products

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
)

type Storage interface {
	GetCatalog(ctx context.Context, department string, category string) ([]*model.CatalogProduct, error)
	GetProduct(ctx context.Context, sku string, attr string) (*model.VariationProduct, error)
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string, category *model.Category) error
	RemoveCategory(ctx context.Context, sku string, category *model.Category) error
}

type Service interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetCatalog(ctx context.Context, department string, category string) ([]*model.CatalogProduct, error)
	GetProduct(ctx context.Context, sku string, attr string) (*model.VariationProduct, error)
	GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string, category *model.Category) error
	RemoveCategory(ctx context.Context, sku string, category *model.Category) error
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

func (s *BasicService) GetCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.storage.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *BasicService) GetCatalog(ctx context.Context, department string, category string) ([]*model.CatalogProduct, error) {
	products, err := s.storage.GetCatalog(ctx, department, category)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *BasicService) GetProduct(ctx context.Context, sku string, attr string) (*model.VariationProduct, error) {
	product, err := s.storage.GetProduct(ctx, sku, attr)
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

func (s *BasicService) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
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

func (s *BasicService) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	err := s.storage.AddCategory(ctx, sku, category)
	if err != nil {
		return err
	}
	return nil
}

func (s *BasicService) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	err := s.storage.RemoveCategory(ctx, sku, category)
	if err != nil {
		return err
	}
	return nil
}
