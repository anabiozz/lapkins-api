package api

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkin-project/lapkin-api/common"
	"github.com/anabiozz/logger"
)

// GetProductVariantByID ...
func GetProductVariantByID(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		product, err := env.DB.GetProductByID(r.URL.Query().Get("product_variant_id"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(product)
	})
}
