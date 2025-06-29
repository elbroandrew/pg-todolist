package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pg-todolist/internal/interfaces"
	"pg-todolist/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxTokenTTL      = 30 * 24 * time.Hour
	revokedPrefix    = "revoked:"
	userPrefix       = "user:"
	lastRefreshKey   = "last_refresh"
	taskUpdatesQueue = "task_updates"
	taskDeletesQueue = "task_deletes"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}
var _ interfaces.TokenRepository = (*RedisRepository)(nil) // Проверка реализации

func NewRedisRepository(addr, password string, db int) (*RedisRepository, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Printf("SUCCESS CONNECTION! REDIS")

	return &RedisRepository{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *RedisRepository) RevokeToken(token string, ttl time.Duration) error {
	if ttl < time.Second {
		ttl = time.Second // Минимальный TTL
	}
	if ttl > maxTokenTTL {
		ttl = maxTokenTTL
	}
	// Добавляем токен в Redis с TTL
	return r.client.SetEx(r.ctx, revokedPrefix+token, "1", ttl).Err()
}

func (r *RedisRepository) IsTokenRevoked(token string) (bool, error) {
	
	// Проверяем наличие токена в отозванных
	exists, err := r.client.Exists(r.ctx, "revoked:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}

func (r *RedisRepository) SetLastRefresh(userID uint, t time.Time) error {
	return r.client.Set(r.ctx, fmt.Sprintf("user:%d:last_refresh", userID), t.Unix(), 0).Err()
}

func (r *RedisRepository) GetLastRefresh(userID uint) (time.Time, error) {
	val, err := r.client.Get(r.ctx, fmt.Sprintf("user:%d:last_refresh", userID)).Int64()
	if err == redis.Nil {
		return time.Now(), nil
	}
	return time.Unix(val, 0), err
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisRepository) GetAllTasks(userID uint) ([]models.Task, error) {
	//получаю все ключи задач для пользователя
	ctx, cancel := context.WithTimeout(r.ctx, 500*time.Millisecond)
	defer cancel()

	// Получаем все ключи с помощью SCAN вместо KEYS (более безопасно для production)
	var keys []string
	iter := r.client.Scan(ctx, 0, fmt.Sprintf("task:%d:*", userID), 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan task keys: %w", err)
	}

	// Используем Pipeline для массового получения
	cmds := make(map[string]*redis.StringCmd, len(keys))
	pipe := r.client.Pipeline()
	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, fmt.Errorf("pipeline exec failed: %w", err)
	}

	// Обрабатываем результаты
	var tasks []models.Task
	var errs []error

	for key, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, fmt.Errorf("key %s: %w", key, err))
			continue
		}

		var task models.Task
		if err := json.Unmarshal([]byte(val), &task); err != nil {
			errs = append(errs, fmt.Errorf("unmarshal error for key %s: %w", key, err))
			continue
		}

		// Проверка целостности данных
		if task.ID == 0 || task.UserID != userID {
			errs = append(errs, fmt.Errorf("invalid task data in key %s", key))
			continue
		}

		// Проверяем что задача не удалена и принадлежит пользователю
		if !task.DeletedAt.Valid && task.UserID == userID {
			tasks = append(tasks, task)
		}
	}

	//  Логируем ошибки, если они есть
	if len(errs) > 0 {
		r.logErrors("GetAllTasks errors:", errs)
	}

	return tasks, nil
}

func (r *RedisRepository) logErrors(prefix string, errs []error) {
	// TODO: использовать мой логгер
	for _, err := range errs {
		log.Printf("%s %v", prefix, err)
	}
}

func (r *RedisRepository) SetEx(key string, value string, expiration time.Duration) error {
	return r.client.SetEx(r.ctx, key, value, expiration).Err()
}

func (r *RedisRepository) PushUpdate(data interface{}) (int64, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("marshal error: %w", err)
	}
	return r.client.LPush(r.ctx, taskUpdatesQueue, jsonData).Result()
}

func (r *RedisRepository) PushDelete(data interface{}) (int64, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("marshal error: %w", err)
	}
	return r.client.LPush(r.ctx, "task_deletes", jsonData).Result()
}
