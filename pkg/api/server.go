package api

import (
	"context"
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
		Handler:      router,
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

	router.Path("/api/v1/products").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetProductsEndpoint(svc),
		decodeGetProductRequest,
		encodeGetProductResponse,
		opts...,
	))

	router.Path("/api/v1/product").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeUpdateProductEndpoint(svc),
		decodeUpdateProductRequest,
		encodeUpdateProductResponse,
		opts...,
	))

	router.Path("/api/v1/products").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetProductsEndpoint(svc),
		decodeGetProductsRequest,
		encodeGetProductsResponse,
		opts...,
	))

	router.Path("/api/v1/categories").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetCategoriesEndpoint(svc),
		decodeGetCategoriesRequest,
		encodeGetCategoriesResponse,
		opts...,
	))

	router.Path("/api/v1/attribute").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeAddAttributeEndpoint(svc),
		decodeAddAttributeRequest,
		encodeAddAttributeResponse,
		opts...,
	))

	router.Path("/api/v1/attribute").Methods(http.MethodDelete).Handler(kithttp.NewServer(
		makeRemoveAttributeEndpoint(svc),
		decodeRemoveAttributeRequest,
		encodeRemoveAttributeResponse,
		opts...,
	))

	router.Path("/api/v1/category").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeAddCategoryEndpoint(svc),
		decodeAddCategoryRequest,
		encodeAddCategoryResponse,
		opts...,
	))

	router.Path("/api/v1/category").Methods(http.MethodDelete).Handler(kithttp.NewServer(
		makeRemoveCategoryEndpoint(svc),
		decodeRemoveCategoryRequest,
		encodeRemoveCategoryResponse,
		opts...,
	))

	router.Path("/api/v1/card/summary").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetHeaderCartInfo(svc),
		decodeGetHeaderCartInfoRequest,
		encodeGetHeaderCartInfoResponse,
		opts...,
	))

	router.Path("/api/v1/card").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetCart(svc),
		decodeGetCartRequest,
		encodeGetCartResponse,
		opts...,
	))

	router.Path("/api/v1/card/inc").Methods(http.MethodPut).Handler(kithttp.NewServer(
		makeIncreaseProductQty(svc),
		decodeIncreaseProductQtyRequest,
		encodeIncreaseProductQtyResponse,
		opts...,
	))

	router.Path("/api/v1/card/dec").Methods(http.MethodPut).Handler(kithttp.NewServer(
		makeDecreaseProductQty(svc),
		decodeDecreaseProductQtyRequest,
		encodeDecreaseProductQtyResponse,
		opts...,
	))

	router.Path("/api/v1/card/product").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeAddProductEndpoint(svc),
		decodeAddProductRequest,
		encodeAddProductResponse,
		opts...,
	))

	router.Path("/api/v1/card/product").Methods(http.MethodDelete).Handler(kithttp.NewServer(
		makeRemoveProduct(svc),
		decodeRemoveProductRequest,
		encodeRemoveProductResponse,
		opts...,
	))

	router.Path("/api/v1/order").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeAddOrderEndpoint(svc),
		decodeAddOrderRequest,
		encodeAddOrderResponse,
		opts...,
	))

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
