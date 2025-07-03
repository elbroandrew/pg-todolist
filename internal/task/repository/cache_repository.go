package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"pg-todolist/internal/models"
	"pg-todolist/pkg/app_errors"
	"time"

	"github.com/redis/go-redis/v9"
)



type cacheRepository struct {
	client *redis.Client
}

func NewTaskCache(client *redis.Client) CacheRepository {
	return &cacheRepository{client: client}
}

func (c *cacheRepository) SetUserTasks(userID uint, tasks []models.Task) *app_errors.AppError {
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		return app_errors.ErrInternalServer
	}

	err = c.client.Set(context.Background(), 
		fmt.Sprintf("user_tasks:%d", userID), 
		jsonData, 
		24*time.Hour,
	).Err()

	if err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

func (c *cacheRepository) GetUserTasks(userID uint) ([]models.Task, *app_errors.AppError) {
	data, err := c.client.Get(context.Background(), 
		fmt.Sprintf("user_tasks:%d", userID),
	).Bytes()

	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, app_errors.ErrDBError
	}

	var tasks []models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, app_errors.ErrInternalServer
	}
	return tasks, nil
}

func (c *cacheRepository) InvalidateUserTasks(userID uint) *app_errors.AppError {
	err := c.client.Del(context.Background(), 
		fmt.Sprintf("user_tasks:%d", userID),
	).Err()

	if err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

func (c *cacheRepository) MarkCompleted(taskID uint) *app_errors.AppError {
	err := c.client.Set(context.Background(), 
		fmt.Sprintf("completed_tasks:%d", taskID), 
		"1", 
		24*time.Hour,
	).Err()

	if err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

func (c *cacheRepository) MarkDeleted(taskID uint) *app_errors.AppError {
	err := c.client.Set(context.Background(), 
		fmt.Sprintf("deleted_tasks:%d", taskID), 
		"1", 
		24*time.Hour,
	).Err()

	if err != nil {
		return app_errors.ErrDBError
	}
	return nil
}