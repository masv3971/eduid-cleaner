package storage

import (
	"context"

	"eduid-cleaner/pkg/model"
	"github.com/go-redis/redis/v8"
)

type KV interface {
	AddToCounter(ctx context.Context, key string) error
	GetCounter(ctx context.Context, key string) (int, error)
	SetStatus(ctx context.Context, nodeName string, statusMSG interface{}) error
	GetAllStatus(ctx context.Context) (map[string]string, error)
}

type Client struct {
	redis *redis.Client
}

func New(config *model.Cfg) (*Client, error) {
	c := &Client{
		redis: redis.NewClient(&redis.Options{
			Addr: config.Storage.Redis.Addr,
			DB:   config.Storage.Redis.DB,
		}),
	}

	return c, nil
}
