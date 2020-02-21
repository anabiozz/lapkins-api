package cart

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkins-api/common"
	"github.com/anabiozz/lapkins-api/models"
	"github.com/anabiozz/logger"
)

// RemoveProduct ...
func RemoveProduct(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		variation := &models.Variation{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(variation)
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		err = env.DB.RemoveProduct(r.URL.Query().Get("cart_session"), variation)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(nil)
	})
}
