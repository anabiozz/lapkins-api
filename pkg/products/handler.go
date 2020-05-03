package products

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/time/rate"
)

type HandlerConfig struct {
	Svc         Service
	Logger      log.Logger
	EnableAuth  bool
	RateLimiter *rate.Limiter
	ReqMetrics  metrics.Histogram
}

// newHandler creates a new HTTP handler serving service endpoints.
func MakeHandler(cfg HandlerConfig) http.Handler {
	var svc Service = &instrumentingMiddleware{
		next:        &loggingMiddleware{Logger: cfg.Logger, next: cfg.Svc},
		reqDuration: cfg.ReqMetrics,
	}

	if cfg.EnableAuth {
		svc = &authMiddleware{
			logger: cfg.Logger,
			next:   svc,
		}
	}

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(makeErrorHandler(cfg.Logger)),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(populateRequestIDIntoContext),
	}

	router := mux.NewRouter()

	getCatalogEndpoint := makeGetCatalogEndpoint(svc)
	getCatalogEndpoint = applyMiddlewares(getCatalogEndpoint, cfg)
	router.Path("/get-catalog").Methods(http.MethodGet).Handler(kithttp.NewServer(
		getCatalogEndpoint,
		decodeGetCatalogRequest,
		encodeGetCatalogResponse,
		opts...,
	))

	getProductEndpoint := makeGetProductEndpoint(svc)
	getProductEndpoint = applyMiddlewares(getProductEndpoint, cfg)
	router.Path("/get-product").Methods(http.MethodGet).Handler(kithttp.NewServer(
		getProductEndpoint,
		decodeGetProductRequest,
		encodeGetProductResponse,
		opts...,
	))

	getCategoryEndpoint := makeGetCategoryEndpoint(svc)
	getCategoryEndpoint = applyMiddlewares(getCategoryEndpoint, cfg)
	router.Path("/get-category").Methods(http.MethodGet).Handler(kithttp.NewServer(
		getCategoryEndpoint,
		decodeGetCategoryRequest,
		encodeGetCategoryResponse,
		opts...,
	))

	getProductsByCategoryEndpoint := makeGetProductsByCategoryEndpoint(svc)
	getProductsByCategoryEndpoint = applyMiddlewares(getProductsByCategoryEndpoint, cfg)
	router.Path("/get-products-by-category").Methods(http.MethodGet).Handler(kithttp.NewServer(
		getProductsByCategoryEndpoint,
		decodeGetProductsByCategoryRequest,
		encodeGetProductsByCategoryResponse,
		opts...,
	))

	addAttributeEndpoint := makeAddAttributeEndpoint(svc)
	addAttributeEndpoint = applyMiddlewares(addAttributeEndpoint, cfg)
	router.Path("/add-attribute").Methods(http.MethodPost).Handler(kithttp.NewServer(
		addAttributeEndpoint,
		decodeAddAttributeRequest,
		encodeAddAttributeResponse,
		opts...,
	))

	removeAttributeEndpoint := makeRemoveAttributeEndpoint(svc)
	removeAttributeEndpoint = applyMiddlewares(removeAttributeEndpoint, cfg)
	router.Path("/remove-attribute").Methods(http.MethodPost).Handler(kithttp.NewServer(
		removeAttributeEndpoint,
		decodeRemoveAttributeRequest,
		encodeRemoveAttributeResponse,
		opts...,
	))

	addCategoryEndpoint := makeAddCategoryEndpoint(svc)
	addCategoryEndpoint = applyMiddlewares(addCategoryEndpoint, cfg)
	router.Path("/add-category").Methods(http.MethodPost).Handler(kithttp.NewServer(
		addCategoryEndpoint,
		decodeAddCategoryRequest,
		encodeAddCategoryResponse,
		opts...,
	))

	removeCategoryEndpoint := makeRemoveCategoryEndpoint(svc)
	removeCategoryEndpoint = applyMiddlewares(removeCategoryEndpoint, cfg)
	router.Path("/remove-category").Methods(http.MethodPost).Handler(kithttp.NewServer(
		removeCategoryEndpoint,
		decodeRemoveCategoryRequest,
		encodeRemoveCategoryResponse,
		opts...,
	))

	return router
}

func applyMiddlewares(e endpoint.Endpoint, cfg HandlerConfig) endpoint.Endpoint {
	e = ratelimit.NewErroringLimiter(cfg.RateLimiter)(e)
	return e
}

type requestIDContextKey struct{}

func getRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDContextKey{}).(string)
	return id, ok
}

func makeErrorHandler(logger log.Logger) kittransport.ErrorHandler {
	return kittransport.ErrorHandlerFunc(func(ctx context.Context, err error) {
		requestID, _ := getRequestID(ctx)
		level.Error(logger).Log("requestID", requestID, "err", err)
	})
}

func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, id)
}

func populateRequestIDIntoContext(ctx context.Context, _ *http.Request) context.Context {
	return withRequestID(ctx, uuid.NewV4().String())
}
