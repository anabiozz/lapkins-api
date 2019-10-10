package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anabiozz/lapkin-project/lapkin-api/common"
	"github.com/anabiozz/logger"
)

// GetProducts ...
func GetProducts(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println(r.URL.Query().Get("products_type"))

		products, err := env.DB.GetProducts(r.URL.Query().Get("products_type"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		log.Println(products)

		json.NewEncoder(w).Encode(products)
	})
}
