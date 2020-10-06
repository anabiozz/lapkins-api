package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/products"
	"github.com/anabiozz/lapkins-api/pkg/storage/mongo"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"

	_ "github.com/lib/pq"
)

const (
	// URL ...
	URL = "localhost:8084"
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

	//cartsSvc, err := cart.NewService(cart.ServiceConfig{
	//	Logger:  logger,
	//	Storage: storage,
	//})
	//if err != nil {
	//	level.Error(logger).Log("msg", "failed to initialize cart service", "err", err)
	//	os.Exit(1)
	//}
	//
	//cartsHandler := cart.MakeHandler(cart.HandlerConfig{
	//	Svc:         cartsSvc,
	//	Logger:      logger,
	//	EnableAuth:  false,
	//	RateLimiter: rate.NewLimiter(rate.Every(100*time.Nanosecond), 100),
	//	ReqMetrics:  reqMetrics,
	//})

	productsSvc, err := products.NewService(products.ServiceConfig{
		Logger:  logger,
		Storage: storage,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize products service", "err", err)
		os.Exit(1)
	}

	productsHandler := products.MakeHandler(products.HandlerConfig{
		Svc:         productsSvc,
		Logger:      logger,
		EnableAuth:  false,
		RateLimiter: rate.NewLimiter(rate.Every(100*time.Nanosecond), 100),
		ReqMetrics:  reqMetrics,
	})

	//usersSvc, err := users.NewService(users.ServiceConfig{
	//	Logger:  logger,
	//	Storage: storage,
	//})
	//if err != nil {
	//	level.Error(logger).Log("msg", "failed to initialize users service", "err", err)
	//	os.Exit(1)
	//}

	//usersHandler := users.MakeHandler(users.HandlerConfig{
	//	Svc:         usersSvc,
	//	Logger:      logger,
	//	RateLimiter: rate.NewLimiter(rate.Every(100*time.Nanosecond), 100),
	//	ReqMetrics:  reqMetrics,
	//})

	//ordersSvc, err := orders.NewService(orders.ServiceConfig{
	//	Logger:  logger,
	//	Storage: storage,
	//})
	//if err != nil {
	//	level.Error(logger).Log("msg", "failed to initialize users service", "err", err)
	//	os.Exit(1)
	//}

	//ordersHandler := orders.MakeHandler(orders.HandlerConfig{
	//	Svc:         ordersSvc,
	//	Logger:      logger,
	//	RateLimiter: rate.NewLimiter(rate.Every(100*time.Nanosecond), 100),
	//	ReqMetrics:  reqMetrics,
	//})

	m := http.NewServeMux()
	//m.Handle("/api/v1/cart/", http.StripPrefix("/api/v1/cart", cors(cartsHandler)))
	m.Handle("/api/v1/", cors(productsHandler))
	//m.Handle("/api/v1/users/", http.StripPrefix("/api/v1/users", cors(usersHandler)))
	//m.Handle("/api/v1/orders/", http.StripPrefix("/api/v1/orders", cors(ordersHandler)))
	//m.Handle("/metrics", promhttp.Handler())

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	m.Handle("/liveness", okHandler)
	m.Handle("/readiness", okHandler)

	srv := http.Server{
		Addr:    URL,
		Handler: m,
	}

	go func() {
		level.Info(logger).Log("msg", "starting server")

		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				level.Info(logger).Log("msg", "server closed")
				return
			}
			level.Error(logger).Log("err", err)
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
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
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
