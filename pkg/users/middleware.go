package users

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

func (mw *loggingMiddleware) Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	begin := time.Now()
	uInfo, err := mw.next.Register(ctx, input)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "Register",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return uInfo, err
}

func (mw *loggingMiddleware) Login(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	begin := time.Now()
	uInfo, err := mw.next.Login(ctx, input)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "Login",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return uInfo, err
}

func (mw *loggingMiddleware) RefreshToken(ctx context.Context, token string) (*model.UserOutput, error) {
	begin := time.Now()
	uInfo, err := mw.next.RefreshToken(ctx, token)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "RefreshToken",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return uInfo, err
}

func (mw *loggingMiddleware) GetUsers(ctx context.Context) ([]*model.User, error) {
	begin := time.Now()
	uInfo, err := mw.next.GetUsers(ctx)
	duration := time.Since(begin).Seconds()
	lvl := level.Debug
	if err != nil || duration > 1 {
		lvl = level.Error
	}
	requestID, _ := getRequestID(ctx)
	lvl(mw).Log(
		"method", "GetUsers",
		"requestID", requestID,
		"err", err,
		"took", duration,
	)
	return uInfo, err
}

// instrumentingMiddleware wraps Service and records request count and duration.
type instrumentingMiddleware struct {
	next        Service
	reqDuration metrics.Histogram
}

func (mw *instrumentingMiddleware) Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.Register(ctx, input)
	labels := []string{"method", "Register", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *instrumentingMiddleware) Login(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.Login(ctx, input)
	labels := []string{"method", "Login", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *instrumentingMiddleware) RefreshToken(ctx context.Context, token string) (*model.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.RefreshToken(ctx, token)
	labels := []string{"method", "RefreshToken", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *instrumentingMiddleware) GetUsers(ctx context.Context) ([]*model.User, error) {
	begin := time.Now()
	resp, err := mw.next.GetUsers(ctx)
	labels := []string{"method", "GetUsers", "error", strconv.FormatBool(err != nil)}
	mw.reqDuration.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

type authMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw *authMiddleware) Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	return mw.next.Register(ctx, input)
}

func (mw *authMiddleware) Login(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	return mw.next.Login(ctx, input)
}

func (mw *authMiddleware) RefreshToken(ctx context.Context, token string) (*model.UserOutput, error) {
	return mw.next.RefreshToken(ctx, token)
}

func (mw *authMiddleware) GetUsers(ctx context.Context) ([]*model.User, error) {
	return mw.next.GetUsers(ctx)
}
