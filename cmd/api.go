package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/carts"
	"github.com/anabiozz/lapkins-api/pkg/products"
	"github.com/anabiozz/lapkins-api/pkg/storage/mongo"
	"github.com/anabiozz/lapkins-api/pkg/users"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"

	_ "github.com/lib/pq"
)

const (
	// URL ...
	URL = "localhost:8081"
)

const (
	metricsNamespace = "lapkins"
	metricsSubsystem = "api"
)

func main() {
	ctx := context.Background()

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	storage, err := mongo.New(mongo.Config{
		Logger: logger,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize mongo storage", "err", err)
		os.Exit(1)
	}

	reqMetrics := kitprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Subsystem: metricsSubsystem,
		Name:      "request_duration",
		Help:      "Request duration",
		Buckets:   []float64{0.001, 0.002, 0.003, 0.004, 0.005, 0.010, 0.015, 0.030, 0.050, 0.100, 0.3, 0.5, 0.8, 1.2, 2},
	}, []string{"method", "error"})

	cartsSvc, err := carts.NewService(carts.ServiceConfig{
		Logger:  logger,
		Storage: storage,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize carts service", "err", err)
		os.Exit(1)
	}

	cartsHandler := carts.MakeHandler(carts.HandlerConfig{
		Svc:         cartsSvc,
		Logger:      logger,
		EnableAuth:  false,
		RateLimiter: rate.NewLimiter(rate.Every(100*time.Nanosecond), 100),
		ReqMetrics:  reqMetrics,
	})

	productsSvc := products.NewService(logger)
	productHandler := products.MakeHandler(ctx, productsSvc, storage)

	usersSvc, err := users.NewService(users.ServiceConfig{
		Logger:  logger,
		Storage: storage,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize users service", "err", err)
		os.Exit(1)
	}

	usersHandler := users.MakeHandler(users.HandlerConfig{
		Svc:         usersSvc,
		Logger:      logger,
		RateLimiter: rate.NewLimiter(rate.Every(100*time.Nanosecond), 100),
		ReqMetrics:  reqMetrics,
	})

	m := http.NewServeMux()
	m.Handle("/api/v1/carts/", http.StripPrefix("/api/v1/carts", cors(cartsHandler)))
	m.Handle("/api/v1/products/", http.StripPrefix("/api/v1", productHandler))
	m.Handle("/api/v1/users/", http.StripPrefix("/api/v1/users", cors(usersHandler)))
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
		level.Info(logger).Log("msg", "starting server")

		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				level.Info(logger).Log("msg", "server closed")
				return
			}
			level.Error(logger).Log("err", "failed to serve service")
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	level.Info(logger).Log("msg", "got a signal, server shutdown", "sig", <-sigCh)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		level.Error(logger).Log("err", "failed to shutdown server")
		os.Exit(1)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X_Requested-With, Accept, Z-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
