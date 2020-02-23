package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anabiozz/lapkins-api/common"
	"github.com/anabiozz/logger"
)

// IncreaseProductQuantity ...
func IncreaseProductQuantity(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		variationID, err := strconv.Atoi(r.URL.Query().Get("variation_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}
		sizeOptionID, err := strconv.Atoi(r.URL.Query().Get("size_option_id"))
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		err = env.DB.IncreaseProductQuantity(variationID, r.URL.Query().Get("cart_session"), sizeOptionID)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(true)
	})
}
