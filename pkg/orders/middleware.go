package orders

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

func (mw *loggingMiddleware) AddOrder(ctx context.Context, order *model.Order) error {
	begin := time.Now()
	err := mw.next.AddOrder(ctx, order)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "AddOrder",
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

func (mw *instrumentingMiddleware) AddOrder(ctx context.Context, order *model.Order) error {
	begin := time.Now()
	err := mw.next.AddOrder(ctx, order)
	labels := []string{"method", "AddOrder", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

type authMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw *authMiddleware) AddOrder(ctx context.Context, order *model.Order) error {
	return mw.next.AddOrder(ctx, order)
}
