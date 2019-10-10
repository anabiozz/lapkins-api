package cart

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/anabiozz/lapkin-project/lapkin-api/common"
	"github.com/anabiozz/logger"
)

// AddProduct ...
func AddProduct(env *common.Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		respBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		var raw map[string]interface{}
		json.Unmarshal(respBody, &raw)

		marshaledVariant, err := json.Marshal(raw["product"])
		if err != nil {
			logger.Info(err)
			json.NewEncoder(w).Encode(logger.Return(err))
			return
		}

		cartSession, err := env.DB.AddProduct(marshaledVariant)
		if err != nil {
			logger.Info(err)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
		}

		json.NewEncoder(w).Encode(cartSession)
	})
}
