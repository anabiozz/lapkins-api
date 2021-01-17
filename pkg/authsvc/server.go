package erpsvc

import (
	"context"
	"github.com/gorilla/handlers"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-kit/kit/log"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ServerConfig is a server configuration.
type ServerConfig struct {
	Logger          log.Logger
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	MetricPrefix    string
	AllowedOrigins  []string
}

// Server is a service server.
type Server struct {
	cfg *ServerConfig
	srv *http.Server
}

// NewServer creates a new server.
func NewServer(cfg ServerConfig) (*Server, error) {
	var svc Service
	svc, err := newService(&ServiceConfig{
		Logger: cfg.Logger,
	})
	if err != nil {
		return nil, err
	}

	svc = NewLoggingMiddleware(svc, cfg.Logger)
	svc = NewInstrumentingMiddleware(svc, cfg.MetricPrefix+"_api")

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/api/v1/", makeHandler(svc))

	srv := &http.Server{
		Handler:      handlers.CORS(handlers.AllowedOrigins(cfg.AllowedOrigins))(router),
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	s := &Server{
		cfg: &cfg,
		srv: srv,
	}

	return s, nil
}

// Serve starts the HTTP server.
func (s *Server) Serve() error {
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown stops the HTTP server.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func makeHandler(svc Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	router := mux.NewRouter()

	router.Path("/api/v1/user/register").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeRegisterEndpoint(svc),
		decodeRegisterRequest,
		encodeRegisterResponse,
		opts...,
	))

	router.Path("/api/v1/user/login").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest,
		encodeLoginResponse,
		opts...,
	))

	router.Path("/api/v1/user/refresh-token").Methods(http.MethodPut).Handler(kithttp.NewServer(
		makeRefreshTokenEndpoint(svc),
		decodeRefreshTokenRequest,
		encodeRefreshTokenResponse,
		opts...,
	))

	router.Path("/api/v1/users").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetUsersEndpoint(svc),
		decodeGetUsersRequest,
		encodeGetUsersResponse,
		opts...,
	))

	return router
}
