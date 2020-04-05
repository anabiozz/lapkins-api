package cart

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func MakeHandler(ctx context.Context, s Service, st Storager) *mux.Router {
	router := mux.NewRouter()
	cart := router.PathPrefix("/cart").Subrouter()
	cart.Handle("/add-product", Cors(s.AddProduct(st))).Methods("POST", "OPTIONS")
	cart.Handle("/increase-product-quantity", Cors(s.IncreaseProductQuantity(st))).Methods("POST", "OPTIONS")
	cart.Handle("/decrease-product-quantity", Cors(s.DecreaseProductQuantity(st))).Methods("POST", "OPTIONS")
	cart.Handle("/remove-product", Cors(s.RemoveProduct(st))).Methods("POST", "OPTIONS")
	cart.Handle("/load-cart", Cors(s.LoadCart(st))).Methods("GET", "OPTIONS")
	cart.Handle("/create-order", Cors(s.CreateOrder(st))).Methods("POST", "OPTIONS")
	return cart
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
