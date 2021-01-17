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

func (mw *LoggingMiddleware) Register(ctx context.Context, input *erp.UserInput) (*erp.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.Register(ctx, input)
	if err != nil {
		level.Error(mw.logger).Log("method", "Register", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) Login(ctx context.Context, input *erp.UserInput, tmpUserID string) (*erp.UserOutput, bool, error) {
	begin := time.Now()
	resp, resp2, err := mw.next.Login(ctx, input, tmpUserID)
	if err != nil {
		level.Error(mw.logger).Log("method", "Login", "err", err, "took", time.Since(begin))
	}
	return resp, resp2, err
}

func (mw *LoggingMiddleware) RefreshToken(ctx context.Context, token string) (*erp.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.RefreshToken(ctx, token)
	if err != nil {
		level.Error(mw.logger).Log("method", "RefreshToken", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) GetUsers(ctx context.Context) ([]*erp.User, error) {
	begin := time.Now()
	resp, err := mw.next.GetUsers(ctx)
	if err != nil {
		level.Error(mw.logger).Log("method", "GetUsers", "err", err, "took", time.Since(begin))
	}
	return resp, err
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

func (mw *InstrumentingMiddleware) Register(ctx context.Context, input *erp.UserInput) (*erp.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.Register(ctx, input)
	labels := []string{"method", "Register", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *InstrumentingMiddleware) Login(ctx context.Context, input *erp.UserInput, tmpUserID string) (*erp.UserOutput, bool, error) {
	begin := time.Now()
	resp, resp2, err := mw.next.Login(ctx, input, tmpUserID)
	labels := []string{"method", "Login", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, resp2, err
}

func (mw *InstrumentingMiddleware) RefreshToken(ctx context.Context, token string) (*erp.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.RefreshToken(ctx, token)
	labels := []string{"method", "RefreshToken", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *InstrumentingMiddleware) GetUsers(ctx context.Context) ([]*erp.User, error) {
	begin := time.Now()
	resp, err := mw.next.GetUsers(ctx)
	labels := []string{"method", "GetUsers", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}
