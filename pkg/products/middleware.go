package products

import (
	"context"
	"strconv"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
)

type loggingMiddleware struct {
	log.Logger
	next Service
}

func (mw *loggingMiddleware) GetCatalog(ctx context.Context, category string) ([]*model.CatalogProduct, error) {
	begin := time.Now()
	products, err := mw.next.GetCatalog(ctx, category)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "GetCatalog",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return products, err
}

func (mw *loggingMiddleware) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	begin := time.Now()
	product, err := mw.next.GetProduct(ctx, sku)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "GetProduct",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return product, err
}

func (mw *loggingMiddleware) GetCategory(ctx context.Context, category string) ([]*model.Category, error) {
	begin := time.Now()
	categories, err := mw.next.GetCategory(ctx, category)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "GetCategory",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return categories, err
}

func (mw *loggingMiddleware) GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error) {
	begin := time.Now()
	products, err := mw.next.GetProductsByCategory(ctx, category)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "GetProductsByCategory",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return products, err
}

func (mw *loggingMiddleware) AddAttribute(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.AddAttribute(ctx, sku)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "AddAttribute",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

func (mw *loggingMiddleware) RemoveAttribute(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.RemoveAttribute(ctx, sku)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "RemoveAttribute",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

func (mw *loggingMiddleware) AddCategory(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "AddCategory",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

func (mw *loggingMiddleware) RemoveCategory(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "RemoveCategory",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

// instrumentingMiddleware wraps Service and records request count and duration.
type instrumentingMiddleware struct {
	next        Service
	reqDuration metrics.Histogram
}

func (mw *instrumentingMiddleware) GetCatalog(ctx context.Context, category string) ([]*model.CatalogProduct, error) {
	begin := time.Now()
	products, err := mw.next.GetCatalog(ctx, category)
	labels := []string{"method", "GetCatalog", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return products, err
}

func (mw *instrumentingMiddleware) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	begin := time.Now()
	product, err := mw.next.GetProduct(ctx, sku)
	labels := []string{"method", "GetProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return product, err
}

func (mw *instrumentingMiddleware) GetCategory(ctx context.Context, category string) ([]*model.Category, error) {
	begin := time.Now()
	categories, err := mw.next.GetCategory(ctx, category)
	labels := []string{"method", "GetCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return categories, err
}

func (mw *instrumentingMiddleware) GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error) {
	begin := time.Now()
	products, err := mw.next.GetProductsByCategory(ctx, category)
	labels := []string{"method", "GetProductsByCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return products, err
}

func (mw *instrumentingMiddleware) AddAttribute(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.AddAttribute(ctx, sku)
	labels := []string{"method", "AddAttribute", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) RemoveAttribute(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.RemoveAttribute(ctx, sku)
	labels := []string{"method", "RemoveAttribute", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) AddCategory(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku)
	labels := []string{"method", "AddCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) RemoveCategory(ctx context.Context, sku string) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku)
	labels := []string{"method", "RemoveCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

type authMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw *authMiddleware) GetCatalog(ctx context.Context, category string) ([]*model.CatalogProduct, error) {
	return mw.next.GetCatalog(ctx, category)
}

func (mw *authMiddleware) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	return mw.next.GetProduct(ctx, sku)
}

func (mw *authMiddleware) GetCategory(ctx context.Context, category string) ([]*model.Category, error) {
	return mw.next.GetCategory(ctx, category)
}

func (mw *authMiddleware) GetProductsByCategory(ctx context.Context, category string) ([]*model.SKUProduct, error) {
	return mw.next.GetProductsByCategory(ctx, category)
}

func (mw *authMiddleware) AddAttribute(ctx context.Context, sku string) error {
	return mw.next.AddAttribute(ctx, sku)
}

func (mw *authMiddleware) RemoveAttribute(ctx context.Context, sku string) error {
	return mw.next.RemoveAttribute(ctx, sku)
}

func (mw *authMiddleware) AddCategory(ctx context.Context, sku string) error {
	return mw.next.AddCategory(ctx, sku)
}

func (mw *authMiddleware) RemoveCategory(ctx context.Context, sku string) error {
	return mw.next.RemoveCategory(ctx, sku)
}
