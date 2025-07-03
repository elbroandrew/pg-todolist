package database

import (
	"context"
	"time"
	"pg-todolist/pkg/config"

	"github.com/redis/go-redis/v9"
)

type RedisWrapper struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.RedisConfig) (RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
		// PoolSize: cfg.PoolSize,
	})

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisWrapper{client: client}, nil
}

func (r *RedisWrapper) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisWrapper) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisWrapper) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

func (r *RedisWrapper) Close() error {
	return r.client.Close()
}