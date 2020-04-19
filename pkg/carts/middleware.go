package carts

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

func (mw *loggingMiddleware) AddProduct(ctx context.Context, sku string, user *model.CartUser) (*model.CartUser, error) {
	begin := time.Now()
	user, err := mw.next.AddProduct(ctx, sku, user)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "AddProduct",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return user, err
}

func (mw *loggingMiddleware) CreateOrder(ctx context.Context) error {
	begin := time.Now()
	err := mw.next.CreateOrder(ctx)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "CreateOrder",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

func (mw *loggingMiddleware) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.DecreaseProductQuantity(ctx, userID, sku)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "DecreaseProductQuantity",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

func (mw *loggingMiddleware) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
	begin := time.Now()
	cartProducts, err := mw.next.LoadCart(ctx, userID)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "LoadCart",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return cartProducts, err
}

func (mw *loggingMiddleware) RemoveProduct(ctx context.Context) error {
	begin := time.Now()
	err := mw.next.RemoveProduct(ctx)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "RemoveProduct",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

func (mw *loggingMiddleware) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.IncreaseProductQuantity(ctx, userID, sku)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "IncreaseProductQuantity",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return err
}

func (mw *loggingMiddleware) GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error) {
	begin := time.Now()
	info, err := mw.next.GetHeaderCartInfo(ctx, userID)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "GetHeaderCartInfo",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return info, err
}

// instrumentingMiddleware wraps Service and records request count and duration.
type instrumentingMiddleware struct {
	next        Service
	reqDuration metrics.Histogram
}

func (mw *instrumentingMiddleware) AddProduct(ctx context.Context, sku string, user *model.CartUser) (*model.CartUser, error) {
	begin := time.Now()
	user, err := mw.next.AddProduct(ctx, sku, user)
	labels := []string{"method", "AddProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return user, err
}

func (mw *instrumentingMiddleware) CreateOrder(ctx context.Context) error {
	begin := time.Now()
	err := mw.next.CreateOrder(ctx)
	labels := []string{"method", "CreateOrder", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.DecreaseProductQuantity(ctx, userID, sku)
	labels := []string{"method", "DecreaseProductQuantity", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	begin := time.Now()
	err := mw.next.IncreaseProductQuantity(ctx, userID, sku)
	labels := []string{"method", "IncreaseProductQuantity", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
	begin := time.Now()
	cartProducts, err := mw.next.LoadCart(ctx, userID)
	labels := []string{"method", "LoadCart", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return cartProducts, err
}

func (mw *instrumentingMiddleware) RemoveProduct(ctx context.Context) error {
	begin := time.Now()
	err := mw.next.RemoveProduct(ctx)
	labels := []string{"method", "RemoveProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *instrumentingMiddleware) GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error) {
	begin := time.Now()
	info, err := mw.next.GetHeaderCartInfo(ctx, userID)
	labels := []string{"method", "GetHeaderCartInfo", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return info, err
}

type authMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw *authMiddleware) AddProduct(ctx context.Context, sku string, user *model.CartUser) (*model.CartUser, error) {
	return mw.next.AddProduct(ctx, sku, user)
}

func (mw *authMiddleware) CreateOrder(ctx context.Context) error {
	return mw.next.CreateOrder(ctx)
}

func (mw *authMiddleware) DecreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	return mw.next.DecreaseProductQuantity(ctx, userID, sku)
}

func (mw *authMiddleware) IncreaseProductQuantity(ctx context.Context, userID string, sku string) error {
	return mw.next.IncreaseProductQuantity(ctx, userID, sku)
}

func (mw *authMiddleware) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
	return mw.next.LoadCart(ctx, userID)
}

func (mw *authMiddleware) RemoveProduct(ctx context.Context) error {
	return mw.next.RemoveProduct(ctx)
}

func (mw *authMiddleware) GetHeaderCartInfo(ctx context.Context, userID string) (*model.HeaderCartInfo, error) {
	return mw.next.GetHeaderCartInfo(ctx, userID)
}
