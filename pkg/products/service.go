package products

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Storager interface {
	GetProducts(ctx context.Context, category string) (products []*model.CatalogProduct, err error)
	GetVariation(tx context.Context, sku string) (product *model.DescriptionProduct, err error)
	GetCategories(categoryURL string) (categories []model.Category, err error)
}

type Service interface {
	GetCategories(storager Storager) http.HandlerFunc
	GetProducts(ctx context.Context, storager Storager) http.HandlerFunc
	GetVariation(ctx context.Context, storager Storager) http.HandlerFunc
}

type BasicService struct {
	logger log.Logger
}

func NewService(logger log.Logger) *BasicService {
	return &BasicService{logger: logger}
}

func (s *BasicService) GetCategories(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		categories, err := st.GetCategories(r.URL.Query().Get("category_url"))
		if err != nil {
			level.Error(s.logger).Log("err", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(categories)
	})
}

func (s *BasicService) GetProducts(ctx context.Context, st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		products, err := st.GetProducts(ctx, r.URL.Query().Get("category"))
		if err != nil {
			level.Error(s.logger).Log("err", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(products)
	})
}

func (s *BasicService) GetVariation(ctx context.Context, st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		product, err := st.GetVariation(ctx, r.URL.Query().Get("sku"))
		if err != nil {
			level.Error(s.logger).Log("err", err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(product)
	})
}
