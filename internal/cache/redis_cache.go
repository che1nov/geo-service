package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisCache) Get(key string, result interface{}) error {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

func (c *RedisCache) Set(key string, value interface{}) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}
