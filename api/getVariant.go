package api

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkins-api/common"
	"github.com/anabiozz/logger"
)

// GetVariant ...
func GetVariant(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		product, err := env.DB.GetVariant(r.URL.Query().Get("variant_id"), r.URL.Query().Get("size_option_id"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(product)
	})
}
