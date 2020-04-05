package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/cart"
	"github.com/anabiozz/lapkins-api/pkg/products"
	"github.com/anabiozz/lapkins-api/pkg/storage/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	// URL ...
	URL = "localhost:8081"
)

func main() {
	ctx := context.Background()
	logger := logrus.New()
	cfg := &postgres.Config{Logger: logger}

	st, err := postgres.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	cartSvc := cart.NewService(logger)
	cartHandler := cart.MakeHandler(ctx, cartSvc, st)

	productsSvc := products.NewService(logger)
	productHandler := products.MakeHandler(ctx, productsSvc, st)

	m := http.NewServeMux()
	m.Handle("/api/v1/cart/", http.StripPrefix("/api/v1", cartHandler))
	m.Handle("/api/v1/products/", http.StripPrefix("/api/v1", productHandler))
	m.Handle("/metrics", promhttp.Handler())

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m.Handle("/liveness", okHandler)
	m.Handle("/readiness", okHandler)

	srv := http.Server{
		Addr:    "localhost:8081",
		Handler: m,
	}

	go func() {
		logger.Info("starting server")

		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logger.Info("server closed")

				return
			}

			logger.Error("failed to serve service")

			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	logger.Info("msg", "got a signal, server shutdown", "sig", <-sigCh)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server")
		os.Exit(1)
	}
}
