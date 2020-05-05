package orders

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

	addOrderEndpoint := makeAddOrderEndpoint(svc)
	addOrderEndpoint = applyMiddlewares(addOrderEndpoint, cfg)
	router.Path("/order").Methods(http.MethodPost).Handler(kithttp.NewServer(
		addOrderEndpoint,
		decodeAddOrderRequest,
		encodeAddOrderResponse,
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
