package main

import (
	"net/http"
	"os"
	"time"

	"github.com/anabiozz/lapkin-project/lapkin-api/api"
	"github.com/anabiozz/lapkin-project/lapkin-api/api/cart"
	"github.com/anabiozz/lapkin-project/lapkin-api/common"
	"github.com/anabiozz/lapkin-project/lapkin-api/common/datastore"
	"github.com/anabiozz/lapkin-project/lapkin-api/middleware"
	"github.com/anabiozz/logger"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	// URL ...
	URL = "0.0.0.0:8081"
)

func main() {

	logger.Init(os.Stdout, os.Stdout, os.Stderr, os.Stderr)

	db, err := datastore.NewDatastore(datastore.POSTGRES)
	if err != nil {
		logger.Error(err)
	}
	defer db.CloseDB()

	env := common.Env{DB: db}

	// Create router
	router := mux.NewRouter()

	// API handlers
	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Handle("/get-products", middleware.Cors(api.GetProducts(&env)))
	apiRouter.Handle("/get-variant", middleware.Cors(api.GetVariant(&env)))
	apiRouter.Handle("/get-categories", middleware.Cors(api.GetCategories(&env)))

	cartRouter := apiRouter.PathPrefix("/cart/").Subrouter()
	cartRouter.Handle("/add-product", middleware.Cors(cart.AddProduct(&env)))
	cartRouter.Handle("/change-quantity", middleware.Cors(cart.ChangeQuantity(&env)))
	cartRouter.Handle("/remove-product", middleware.Cors(cart.RemoveProduct(&env)))
	cartRouter.Handle("/get-cart", middleware.Cors(cart.GetCart(&env)))
	cartRouter.Handle("/create-cart-session", middleware.Cors(cart.CreateSession(&env)))

	srv := &http.Server{
		Handler:      router,
		Addr:         URL,
		IdleTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Infof("Server was started on http://%s", URL)
	logger.Fatal(srv.ListenAndServe())
}
