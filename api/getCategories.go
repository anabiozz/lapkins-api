package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anabiozz/lapkin-project/lapkin-api/common"
	"github.com/anabiozz/logger"
)

// GetCategories ...
func GetCategories(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		categories, err := env.DB.GetCategories(r.URL.Query().Get("category_id"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		log.Println(categories)

		json.NewEncoder(w).Encode(categories)
	})
}
