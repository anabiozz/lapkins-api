package api

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkins-api/common"
	"github.com/anabiozz/logger"
)

// GetProducts ...
func GetProducts(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		products, err := env.DB.GetProducts(r.URL.Query().Get("category_url"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(products)
	})
}
