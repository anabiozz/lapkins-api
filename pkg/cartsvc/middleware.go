package erpsvc

import (
	"context"
	"github.com/anabiozz/core/lapkins/pkg/erp"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

func NewLoggingMiddleware(next Service, logger log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
		next:   next,
	}
}

type LoggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw *LoggingMiddleware) AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error) {
	begin := time.Now()
	setTmpUserIDCookie, userID, err := mw.next.AddProductToCard(ctx, sku, userID, isLoggedIn, isTmpUserIDSet)
	if err != nil {
		level.Error(mw.logger).Log("method", "AddProductToCard", "err", err, "took", time.Since(begin))
	}
	return setTmpUserIDCookie, userID, err
}

func (mw *LoggingMiddleware) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.DecreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		level.Error(mw.logger).Log("method", "DecreaseProductQuantity", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.IncreaseProductQuantity(ctx, userID, sku)
	if err != nil {
		level.Error(mw.logger).Log("method", "IncreaseProductQuantity", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) LoadCart(ctx context.Context, userID string) ([]*erp.CartProduct, error) {
	begin := time.Now()
	resp, err := mw.next.LoadCart(ctx, userID)
	if err != nil {
		level.Error(mw.logger).Log("method", "LoadCart", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) RemoveProduct(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.RemoveProduct(ctx, userID, sku)
	if err != nil {
		level.Error(mw.logger).Log("method", "RemoveProduct", "err", err, "took", time.Since(begin))
	}
	return err
}

func NewInstrumentingMiddleware(next Service, prefix string) *InstrumentingMiddleware {
	return &InstrumentingMiddleware{
		next: next,
		reqMetrics: kitprometheus.NewHistogramFrom(
			prometheus.HistogramOpts{
				Name: prefix + "_requests",
				Help: "Requests Info",
				Buckets: []float64{
					0.001, 0.002, 0.003,
					0.004, 0.005, 0.010,
					0.015, 0.030, 0.050,
					0.100, 0.3, 0.5,
					1, 2, 5,
				},
			},
			[]string{"method", "error"},
		),
	}
}

type InstrumentingMiddleware struct {
	next       Service
	reqMetrics metrics.Histogram
}

func (mw *InstrumentingMiddleware) AddProductToCard(ctx context.Context, sku string, userID string, isLoggedIn bool, isTmpUserIDSet bool) (bool, string, error) {
	begin := time.Now()
	setTmpUserIDCookie, userID, err := mw.next.AddProductToCard(ctx, sku, userID, isLoggedIn, isTmpUserIDSet)
	labels := []string{"method", "AddProductToCard", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return setTmpUserIDCookie, userID, err
}

func (mw *InstrumentingMiddleware) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.DecreaseProductQuantity(ctx, userID, sku)
	labels := []string{"method", "DecreaseProductQuantity", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.IncreaseProductQuantity(ctx, userID, sku)
	labels := []string{"method", "IncreaseProductQuantity", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) LoadCart(ctx context.Context, userID string) ([]*erp.CartProduct, error) {
	begin := time.Now()
	resp, err := mw.next.LoadCart(ctx, userID)
	labels := []string{"method", "LoadCart", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *InstrumentingMiddleware) RemoveProduct(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.RemoveProduct(ctx, userID, sku)
	labels := []string{"method", "RemoveProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}
