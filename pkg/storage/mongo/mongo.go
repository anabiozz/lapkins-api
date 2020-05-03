package mongo

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Logger log.Logger
}

type Storage struct {
	db     *mongo.Database
	logger log.Logger
	mu     sync.Mutex
}

func New(cfg Config) (*Storage, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	m := &Storage{
		db:     client.Database("lapkins"),
		logger: cfg.Logger,
	}
	level.Info(cfg.Logger).Log("msg", "mongo was up")
	return m, nil
}
