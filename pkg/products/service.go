package products

import (
	"context"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
)

type Storage interface {
	GetProduct(ctx context.Context, sku string) (*model.Product, error)
	GetProducts(ctx context.Context) ([]*model.Product, error)
	GetCategories(ctx context.Context) ([]*model.Category, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string, category *model.Category) error
	RemoveCategory(ctx context.Context, sku string, category *model.Category) error
	UpdateProduct(ctx context.Context, product *model.Product) error
}

type Service interface {
	GetCategories(ctx context.Context) ([]*model.Category, error)
	GetProducts(ctx context.Context) ([]*model.Product, error)
	GetProduct(ctx context.Context, sku string) (*model.Product, error)
	AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error
	RemoveAttribute(ctx context.Context, sku string, attribute string) error
	AddCategory(ctx context.Context, sku string, category *model.Category) error
	RemoveCategory(ctx context.Context, sku string, category *model.Category) error
	UpdateProduct(ctx context.Context, product *model.Product) error
}

type service struct {
	logger  log.Logger
	storage Storage
}

type ServiceConfig struct {
	Logger  log.Logger
	Storage Storage
}

func NewService(cfg ServiceConfig) (*service, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = log.NewNopLogger()
	}

	if cfg.Storage == nil {
		return nil, errBadRequest("storage must be provided")
	}

	svc := &service{
		logger:  logger,
		storage: cfg.Storage,
	}

	return svc, nil
}

func (s *service) GetCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.storage.GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *service) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	product, err := s.storage.GetProduct(ctx, sku)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *service) GetProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.storage.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *service) UpdateProduct(ctx context.Context, product *model.Product) error {
	err := s.storage.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
	err := s.storage.AddAttribute(ctx, sku, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	err := s.storage.RemoveAttribute(ctx, sku, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	err := s.storage.AddCategory(ctx, sku, category)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	err := s.storage.RemoveCategory(ctx, sku, category)
	if err != nil {
		return err
	}
	return nil
}
