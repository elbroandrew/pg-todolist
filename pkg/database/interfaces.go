package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)


// Database интерфейс для работы с SQL-БД
type Database interface {
	Connect(ctx context.Context) error
    Migrate(models ...interface{}) error
    GetDB() *gorm.DB
	HealthCheck(ctx context.Context) error
    Close() error
}

// RedisClient интерфейс для Redis
type RedisClient interface {
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Get(ctx context.Context, key string) (string, error)
}