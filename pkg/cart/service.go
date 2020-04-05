package cart

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/anabiozz/lapkins-api/pkg/storage"
	"github.com/anabiozz/logger"
	log "github.com/sirupsen/logrus"
)

type Storager interface {
	GetProducts(categoryURL string) (products []storage.Product, err error)
	GetVariation(variationID, size string) (product *storage.Variation, err error)
	GetCategories(categoryURL string) (categories []storage.Category, err error)

	AddProduct(variationID int, сartSession string, sizeOptionID int) (err error)
	IncreaseProductQuantity(variationID int, cartSession string, sizeOptionID int) (err error)
	DecreaseProductQuantity(variationID int, cartSession string, sizeOptionID int) (err error)
	RemoveProduct(variationID int, cartSession string, sizeOptionID int) (err error)
	LoadCart(cartSession string) (cartItems []storage.CartItemResponse, err error)

	CloseDB()
}

type Service interface {
	AddProduct(st Storager) http.HandlerFunc
	CreateOrder(st Storager) http.HandlerFunc
	DecreaseProductQuantity(st Storager) http.HandlerFunc
	IncreaseProductQuantity(st Storager) http.HandlerFunc
	LoadCart(st Storager) http.HandlerFunc
	RemoveProduct(st Storager) http.HandlerFunc
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

func (s *BasicService) AddProduct(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		respBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		cartItem := storage.CartItem{}

		json.Unmarshal(respBody, &cartItem)

		err = st.AddProduct(cartItem.VariationID, cartItem.СartSession, cartItem.SizeOptionID)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(true)
	})
}

func (s *BasicService) CreateOrder(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		respBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		cartItem := storage.CartItem{}

		json.Unmarshal(respBody, &cartItem)

		//err = st.CreateOrder(cartItem.VariationID, cartItem.СartSession, cartItem.SizeOptionID)
		//if err != nil {
		//	logger.Info(err)
		//	w.WriteHeader(http.StatusNotFound)
		//	json.NewEncoder(w).Encode(err)
		//}

		json.NewEncoder(w).Encode(true)
	})
}

func (s *BasicService) DecreaseProductQuantity(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		variationID, err := strconv.Atoi(r.URL.Query().Get("variation_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}
		sizeOptionID, err := strconv.Atoi(r.URL.Query().Get("size_option_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		err = st.DecreaseProductQuantity(variationID, r.URL.Query().Get("cart_session"), sizeOptionID)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(true)
	})
}

func (s *BasicService) IncreaseProductQuantity(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		variationID, err := strconv.Atoi(r.URL.Query().Get("variation_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}
		sizeOptionID, err := strconv.Atoi(r.URL.Query().Get("size_option_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		err = st.IncreaseProductQuantity(variationID, r.URL.Query().Get("cart_session"), sizeOptionID)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(true)
	})
}

func (s *BasicService) LoadCart(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cart, err := st.LoadCart(r.URL.Query().Get("cart_session"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(cart)
	})
}

func (s *BasicService) RemoveProduct(st Storager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		variationID, err := strconv.Atoi(r.URL.Query().Get("variation_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}
		sizeOptionID, err := strconv.Atoi(r.URL.Query().Get("size_option_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		err = st.RemoveProduct(variationID, r.URL.Query().Get("cart_session"), sizeOptionID)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(true)
	})
}
