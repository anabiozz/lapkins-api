package cart

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkin-project/lapkin-api/common"
	"github.com/anabiozz/lapkin-project/lapkin-api/models"
	"github.com/anabiozz/logger"
)

// RemoveProduct ...
func RemoveProduct(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		variant := &models.Variant{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(variant)
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		err = env.DB.RemoveProduct(r.URL.Query().Get("cart_session"), variant)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(nil)
	})
}
