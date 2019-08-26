package api

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkin-project/lapkin/backend/common"
	"github.com/anabiozz/lapkin-project/lapkin/backend/models"
	"github.com/anabiozz/logger"
)

// GetProducts ...
func GetProducts(env *common.Env, paths models.Paths) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		products, err := env.DB.GetProducts(r.URL.Query().Get("products_type"), paths)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(products)
	})
}
