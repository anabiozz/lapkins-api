package api

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkins-api/common"
	"github.com/anabiozz/logger"
)

// GetCategories ...
func GetCategories(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		categories, err := env.DB.GetCategories(r.URL.Query().Get("category_url"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(categories)
	})
}
