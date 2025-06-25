package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func InitRedis(addr, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	//ping
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("FAILED TO CONNECT TO REDIS")
	}
	log.Printf("SUCCESS CONNECTION!")
}

func RevokeToken(token string, ttl time.Duration) error {
	// Добавляем токен в Redis с TTL
	return rdb.Set(ctx, "revoked:"+token, "1", ttl).Err()
}

func IsTokenRevoked(token string) (bool, error) {
	// Проверяем наличие токена в отозванных
	exists, err := rdb.Exists(ctx, "revoked:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func Close() {
	if rdb != nil {
		rdb.Close()
	}
}
