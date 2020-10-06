package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/api"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"
)

const (
	// URL ...
	URL = "localhost:8084"
)

const metricPrefix = "lapkins"

type configuration struct {
	Port            string        `envconfig:"PORT" required:"true" default:"8080"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(1)
	}

	srv, err := api.NewServer(api.ServerConfig{
		Logger:          logger,
		Port:            cfg.Port,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		ShutdownTimeout: cfg.ShutdownTimeout,
		MetricPrefix:    metricPrefix,
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to create api server", "err", err)
		os.Exit(1)
	}

	go func() {
		level.Info(logger).Log("msg", "starting server", "port", cfg.Port)
		if err := srv.Serve(); err != nil {
			level.Error(logger).Log("msg", "server run failure", "err", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	sig := <-c

	level.Info(logger).Log("msg", "received signal, exiting", "signal", sig)

	if err := srv.Shutdown(); err != nil {
		level.Error(logger).Log("msg", "shutdown failure", "err", err)
	}

	level.Info(logger).Log("msg", "goodbye")
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
