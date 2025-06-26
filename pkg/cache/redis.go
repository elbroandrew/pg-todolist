package cache

import (
	"context"
	"fmt"
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
	log.Printf("SUCCESS CONNECTION! REDIS")
}

func RevokeToken(token string, ttl time.Duration) error {
	key := "revoked:" + token
    if ttl < time.Second {
        ttl = time.Second // Минимальный TTL
    }
	// Добавляем токен в Redis с TTL
	return rdb.SetEx(ctx, key, "1", ttl).Err()
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
func SetLastRefresh(userID uint, t time.Time) error {
    return rdb.Set(ctx, fmt.Sprintf("user:%d:last_refresh", userID), t.Unix(), 0).Err()
}

func GetLastRefresh(userID uint) (time.Time, error) {
    val, err := rdb.Get(ctx, fmt.Sprintf("user:%d:last_refresh", userID)).Int64()
	if err == redis.Nil {
        return time.Now(), nil // Первое обновление
    }
    return time.Unix(val, 0), err
}