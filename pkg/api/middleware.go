package api

import (
	"context"
	"strconv"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
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

func (mw *LoggingMiddleware) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	begin := time.Now()
	resp, err := mw.next.GetProduct(ctx, sku)
	if err != nil {
		level.Error(mw.logger).Log("method", "GetProduct", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) GetProducts(ctx context.Context) ([]*model.Product, error) {
	begin := time.Now()
	resp, err := mw.next.GetProducts(ctx)
	if err != nil {
		level.Error(mw.logger).Log("method", "GetProducts", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) GetCategories(ctx context.Context) ([]*model.Category, error) {
	begin := time.Now()
	resp, err := mw.next.GetCategories(ctx)
	if err != nil {
		level.Error(mw.logger).Log("method", "GetCategories", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
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

func (mw *LoggingMiddleware) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku, category)
	if err != nil {
		level.Error(mw.logger).Log("method", "AddCategory", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku, category)
	if err != nil {
		level.Error(mw.logger).Log("method", "RemoveCategory", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) UpdateProduct(ctx context.Context, product *model.Product) error {
	begin := time.Now()
	err := mw.next.UpdateProduct(ctx, product)
	if err != nil {
		level.Error(mw.logger).Log("method", "UpdateProduct", "err", err, "took", time.Since(begin))
	}
	return err
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

func (mw *LoggingMiddleware) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
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

func (mw *LoggingMiddleware) AddOrder(ctx context.Context, order *model.Order) error {
	begin := time.Now()
	err := mw.next.AddOrder(ctx, order)
	if err != nil {
		level.Error(mw.logger).Log("method", "AddOrder", "err", err, "took", time.Since(begin))
	}
	return err
}

func (mw *LoggingMiddleware) Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.Register(ctx, input)
	if err != nil {
		level.Error(mw.logger).Log("method", "Register", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) Login(ctx context.Context, input *model.UserInput, tmpUserID string) (*model.UserOutput, bool, error) {
	begin := time.Now()
	resp, resp2, err := mw.next.Login(ctx, input, tmpUserID)
	if err != nil {
		level.Error(mw.logger).Log("method", "Login", "err", err, "took", time.Since(begin))
	}
	return resp, resp2, err
}

func (mw *LoggingMiddleware) RefreshToken(ctx context.Context, token string) (*model.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.RefreshToken(ctx, token)
	if err != nil {
		level.Error(mw.logger).Log("method", "RefreshToken", "err", err, "took", time.Since(begin))
	}
	return resp, err
}

func (mw *LoggingMiddleware) GetUsers(ctx context.Context) ([]*model.User, error) {
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

func (mw *InstrumentingMiddleware) GetProduct(ctx context.Context, sku string) (*model.Product, error) {
	begin := time.Now()
	product, err := mw.next.GetProduct(ctx, sku)
	labels := []string{"method", "GetProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return product, err
}

func (mw *InstrumentingMiddleware) GetProducts(ctx context.Context) ([]*model.Product, error) {
	begin := time.Now()
	products, err := mw.next.GetProducts(ctx)
	labels := []string{"method", "GetProducts", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return products, err
}

func (mw *InstrumentingMiddleware) GetCategories(ctx context.Context) ([]*model.Category, error) {
	begin := time.Now()
	categories, err := mw.next.GetCategories(ctx)
	labels := []string{"method", "GetCategories", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return categories, err
}

func (mw *InstrumentingMiddleware) AddAttribute(ctx context.Context, sku string, attribute *model.NameValue) error {
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

func (mw *InstrumentingMiddleware) AddCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.AddCategory(ctx, sku, category)
	labels := []string{"method", "AddCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) RemoveCategory(ctx context.Context, sku string, category *model.Category) error {
	begin := time.Now()
	err := mw.next.RemoveCategory(ctx, sku, category)
	labels := []string{"method", "RemoveCategory", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) UpdateProduct(ctx context.Context, product *model.Product) error {
	begin := time.Now()
	err := mw.next.UpdateProduct(ctx, product)
	labels := []string{"method", "UpdateProduct", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
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

func (mw *InstrumentingMiddleware) LoadCart(ctx context.Context, userID string) ([]*model.CartProduct, error) {
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

func (mw *InstrumentingMiddleware) AddOrder(ctx context.Context, order *model.Order) error {
	begin := time.Now()
	err := mw.next.AddOrder(ctx, order)
	labels := []string{"method", "AddOrder", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return err
}

func (mw *InstrumentingMiddleware) Register(ctx context.Context, input *model.UserInput) (*model.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.Register(ctx, input)
	labels := []string{"method", "Register", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *InstrumentingMiddleware) Login(ctx context.Context, input *model.UserInput, tmpUserID string) (*model.UserOutput, bool, error) {
	begin := time.Now()
	resp, resp2, err := mw.next.Login(ctx, input, tmpUserID)
	labels := []string{"method", "Login", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, resp2, err
}

func (mw *InstrumentingMiddleware) RefreshToken(ctx context.Context, token string) (*model.UserOutput, error) {
	begin := time.Now()
	resp, err := mw.next.RefreshToken(ctx, token)
	labels := []string{"method", "RefreshToken", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}

func (mw *InstrumentingMiddleware) GetUsers(ctx context.Context) ([]*model.User, error) {
	begin := time.Now()
	resp, err := mw.next.GetUsers(ctx)
	labels := []string{"method", "GetUsers", "error", strconv.FormatBool(err != nil)}
	mw.reqMetrics.With(labels...).Observe(time.Since(begin).Seconds())
	return resp, err
}