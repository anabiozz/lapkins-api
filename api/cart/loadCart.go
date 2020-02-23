package cart

import (
	"encoding/json"
	"net/http"

	"github.com/anabiozz/lapkins-api/common"
	"github.com/anabiozz/logger"
)

// LoadCart ...
func LoadCart(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cart, err := env.DB.LoadCart(r.URL.Query().Get("cart_session"))
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(cart)
	})
}
