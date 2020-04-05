package products

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func MakeHandler(ctx context.Context, s Service, st Storager) *mux.Router {
	router := mux.NewRouter()
	products := router.PathPrefix("/products").Subrouter()
	products.Handle("/get-products", Cors(s.GetProducts(st))).Methods("GET", "OPTIONS")
	products.Handle("/get-variation", Cors(s.GetVariation(st))).Methods("GET", "OPTIONS")
	products.Handle("/get-categories", Cors(s.GetCategories(st))).Methods("GET", "OPTIONS")
	return products
}

func Cors(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X_Requested-With, Accept, Z-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)

	})
}
