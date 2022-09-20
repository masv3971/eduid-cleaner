package storage

import (
	"context"
)

func (c *Client) AddToCounter(ctx context.Context, key string) error {
	return c.redis.Incr(ctx, key).Err()
}

func (c *Client) GetCounter(ctx context.Context, key string) (int, error) {
	return c.redis.Get(ctx, key).Int()
}

func (c *Client) SetStatus(ctx context.Context, nodeName string, statusMSG interface{}) error {
	return c.redis.HSet(ctx, "status", nodeName, statusMSG).Err()
}

func (c *Client) GetAllStatus(ctx context.Context) (map[string]string, error) {
	return c.redis.HGetAll(ctx, "status").Result()
}
