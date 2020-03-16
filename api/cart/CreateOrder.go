package cart

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/anabiozz/lapkins-api/common"
	"github.com/anabiozz/lapkins-api/models"
	"github.com/anabiozz/logger"
)

// CreateOrder ...
func CreateOrder(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		respBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		cartItem := models.CartItem{}

		json.Unmarshal(respBody, &cartItem)

		err = env.DB.CreateOrder(cartItem.VariationID, cartItem.СartSession, cartItem.SizeOptionID)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(true)
	})
}
