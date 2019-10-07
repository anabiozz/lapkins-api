package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anabiozz/lapkin-project/lapkin-api/api"
	"github.com/anabiozz/lapkin-project/lapkin-api/common"
	"github.com/anabiozz/lapkin-project/lapkin-api/common/datastore"
	"github.com/anabiozz/lapkin-project/lapkin-api/middleware"
	"github.com/anabiozz/lapkin-project/lapkin-api/models"
	"github.com/anabiozz/logger"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	// URL ...
	URL = "127.0.0.1:8081"
)

func main() {

	logger.Init(os.Stdout, os.Stdout, os.Stderr, os.Stderr)

	db, err := datastore.NewDatastore(datastore.POSTGRES)
	if err != nil {
		logger.Error(err)
	}
	defer db.CloseDB()

	env := common.Env{DB: db}

	// Handlers
	paths := models.Paths{}

	// Images directory path
	imagesPath, err := filepath.Abs(os.Getenv("HOME") + "/images")
	if err != nil {
		logger.Error(err)
	}

	paths.FullPath = imagesPath + "/full/"
	paths.PreviewPath = imagesPath + "/preview/"

	// Create router
	router := mux.NewRouter()

	imagesRouter := router.PathPrefix(imagesPath).Subrouter()
	imagesRouter.PathPrefix("/").Handler(http.StripPrefix(imagesPath+"/", http.FileServer(http.Dir(imagesPath))))

	// API handlers
	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.Handle("/get-products", middleware.Cors(api.GetProducts(&env, paths)))
	apiRouter.Handle("/get-variant", middleware.Cors(api.GetVariant(&env)))
	apiRouter.Handle("/get-categories", middleware.Cors(api.GetCategories(&env)))

	srv := &http.Server{
		Handler:      router,
		Addr:         URL,
		IdleTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Infof("Server was started on http://%s", URL)
	logger.Fatal(srv.ListenAndServe())
}
