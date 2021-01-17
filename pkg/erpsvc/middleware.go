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

func (mw *LoggingMiddleware) GetProduct(ctx context.Context, sku string) (*erp.Product, error) {
	begin := time.Now()
	resp, err := mw.next.GetProduct(ctx, sku)
	if err != nil {
		level.Error(mw.logger).Log("method", "GetProduct", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) GetProducts(ctx context.Context) ([]*erp.Product, error) {
	begin := time.Now()
	resp, err := mw.next.GetProducts(ctx)
	if err != nil {
		level.Error(mw.logger).Log("method", "GetProducts", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) GetCategories(ctx context.Context) ([]*erp.Category, error) {
	begin := time.Now()
	resp, err := mw.next.GetCategories(ctx)
	if err != nil {
		level.Error(mw.logger).Log("method", "GetCategories", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) AddAttribute(ctx context.Context, sku string, attribute *erp.NameValue) error {
	begin := time.Now()
	err := mw.next.AddAttribute(ctx, sku, attribute)
	if err != nil {
		level.Error(mw.logger).Log("method", "AddAttribute", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	begin := time.Now()
	err := mw.next.RemoveAttribute(ctx, sku, attribute)
	if err != nil {
		level.Error(mw.logger).Log("method", "RemoveAttribute", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) AddCategory(ctx context.Context, sku string, category *erp.Category) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku, category)
	if err != nil {
		level.Error(mw.logger).Log("method", "AddCategory", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) RemoveCategory(ctx context.Context, sku string, category *erp.Category) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku, category)
	if err != nil {
		level.Error(mw.logger).Log("method", "RemoveCategory", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) UpdateProduct(ctx context.Context, product *erp.Product) error {
	begin := time.Now()
	err := mw.next.UpdateProduct(ctx, product)
	if err != nil {
		level.Error(mw.logger).Log("method", "UpdateProduct", "err", err, "took", time.Since(begin))
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

func (mw *InstrumentingMiddleware) GetProduct(ctx context.Context, sku string) (*erp.Product, error) {
	begin := time.Now()
	product, err := mw.next.GetProduct(ctx, sku)
	labels := []string{"method", "GetProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return product, err
}

func (mw *InstrumentingMiddleware) GetProducts(ctx context.Context) ([]*erp.Product, error) {
	begin := time.Now()
	products, err := mw.next.GetProducts(ctx)
	labels := []string{"method", "GetProducts", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return products, err
}

func (mw *InstrumentingMiddleware) GetCategories(ctx context.Context) ([]*erp.Category, error) {
	begin := time.Now()
	categories, err := mw.next.GetCategories(ctx)
	labels := []string{"method", "GetCategories", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return categories, err
}

func (mw *InstrumentingMiddleware) AddAttribute(ctx context.Context, sku string, attribute *erp.NameValue) error {
	begin := time.Now()
	err := mw.next.AddAttribute(ctx, sku, attribute)
	labels := []string{"method", "AddAttribute", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	begin := time.Now()
	err := mw.next.RemoveAttribute(ctx, sku, attribute)
	labels := []string{"method", "RemoveAttribute", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) AddCategory(ctx context.Context, sku string, category *erp.Category) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku, category)
	labels := []string{"method", "AddCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) RemoveCategory(ctx context.Context, sku string, category *erp.Category) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku, category)
	labels := []string{"method", "RemoveCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) UpdateProduct(ctx context.Context, product *erp.Product) error {
	begin := time.Now()
	err := mw.next.UpdateProduct(ctx, product)
	labels := []string{"method", "UpdateProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}
