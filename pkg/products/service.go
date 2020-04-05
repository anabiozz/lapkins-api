package products

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkins-api/pkg/storage"
	"github.com/anabiozz/logger"
	log "github.com/sirupsen/logrus"
)

type Storager interface {
	GetProducts(categoryURL string) (products []storage.Product, err error)
	GetVariation(variationID, size string) (product *storage.Variation, err error)
	GetCategories(categoryURL string) (categories []storage.Category, err error)
}

type Service interface {
	GetCategories(storager Storager) http.HandlerFunc
	GetProducts(storager Storager) http.HandlerFunc
	GetVariation(storager Storager) http.HandlerFunc
}

type BasicService struct {
	logger *log.Logger
}

func NewService(logger *log.Logger) *BasicService {
	if logger == nil {
		logger = log.New()
	}
	return &BasicService{logger: logger}
}

func (s *BasicService) GetCategories(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		categories, err := st.GetCategories(r.URL.Query().Get("category_url"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(categories)
	})
}

func (s *BasicService) GetProducts(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		products, err := st.GetProducts(r.URL.Query().Get("category_url"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(products)
	})
}

func (s *BasicService) GetVariation(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		product, err := st.GetVariation(r.URL.Query().Get("variation_id"), r.URL.Query().Get("size_option_id"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(product)
	})
}
