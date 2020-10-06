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

func (mw *loggingMiddleware) GetProducts(ctx context.Context) ([]*model.Product, error) {
	begin := time.Now()
	products, err := mw.next.GetProducts(ctx)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "GetProducts",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return products, err
}

func (mw *loggingMiddleware) GetCategories(ctx context.Context) ([]*model.Category, error) {
	begin := time.Now()
	categories, err := mw.next.GetCategories(ctx)
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

func (mw *loggingMiddleware) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
	begin := time.Now()
	err := mw.next.AddAttribute(ctx, sku, attribute)
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

func (mw *loggingMiddleware) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	begin := time.Now()
	err := mw.next.RemoveAttribute(ctx, sku, attribute)
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

func (mw *loggingMiddleware) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku, category)
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

func (mw *loggingMiddleware) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku, category)
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

func (mw *loggingMiddleware) UpdateProduct(ctx context.Context, product *model.Product) error {
	begin := time.Now()
	err := mw.next.UpdateProduct(ctx, product)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "UpdateProduct",
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

func (mw *instrumentingMiddleware) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	begin := time.Now()
	product, err := mw.next.GetProduct(ctx, sku)
	labels := []string{"method", "GetProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return product, err
}

func (mw *instrumentingMiddleware) GetProducts(ctx context.Context) ([]*model.Product, error) {
	begin := time.Now()
	products, err := mw.next.GetProducts(ctx)
	labels := []string{"method", "GetProducts", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return products, err
}

func (mw *instrumentingMiddleware) GetCategories(ctx context.Context) ([]*model.Category, error) {
	begin := time.Now()
	categories, err := mw.next.GetCategories(ctx)
	labels := []string{"method", "GetCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return categories, err
}

func (mw *instrumentingMiddleware) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
	begin := time.Now()
	err := mw.next.AddAttribute(ctx, sku, attribute)
	labels := []string{"method", "AddAttribute", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	begin := time.Now()
	err := mw.next.RemoveAttribute(ctx, sku, attribute)
	labels := []string{"method", "RemoveAttribute", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku, category)
	labels := []string{"method", "AddCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku, category)
	labels := []string{"method", "RemoveCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) UpdateProduct(ctx context.Context, product *model.Product) error {
	begin := time.Now()
	err := mw.next.UpdateProduct(ctx, product)
	labels := []string{"method", "UpdateProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

type authMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw *authMiddleware) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	return mw.next.GetProduct(ctx, sku)
}

func (mw *authMiddleware) GetProducts(ctx context.Context) ([]*model.Product, error) {
	return mw.next.GetProducts(ctx)
}

func (mw *authMiddleware) GetCategories(ctx context.Context) ([]*model.Category, error) {
	return mw.next.GetCategories(ctx)
}

func (mw *authMiddleware) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
	return mw.next.AddAttribute(ctx, sku, attribute)
}

func (mw *authMiddleware) RemoveAttribute(ctx context.Context, sku string, attribute string) error {
	return mw.next.RemoveAttribute(ctx, sku, attribute)
}

func (mw *authMiddleware) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	return mw.next.AddCategory(ctx, sku, category)
}

func (mw *authMiddleware) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	return mw.next.RemoveCategory(ctx, sku, category)
}

func (mw *authMiddleware) UpdateProduct(ctx context.Context, product *model.Product) error {
	return mw.next.UpdateProduct(ctx, product)
}
