package service

import (
	"errors"
	"fmt"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"time"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
	cache    *repository.RedisRepository
}

const (
	taskUpdatesQueue = "task_updates"
	taskDeletesQueue = "task_deletes"
	maxQueueZize     = 100
)

func NewTaskService(taskRepo *repository.TaskRepository, cache *repository.RedisRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo, cache: cache}
}

func (s *TaskService) Create(task *models.Task) error {
	return s.taskRepo.Create(task)
}

func (s *TaskService) GetTaskByUserID(userID uint) ([]models.Task, error) {
	// 1. сделать проверку, есть ли в редис кеше уже таски
	// (logout их записывает и стирает очередь - Поэтому первый запрос
	//отсюда будет загружать с MySQL данные на фронетенд). 
	// брать таски с редиса, если нет их - с БД.
	//причем брать все кроме deleted 

	redisTasks, err := s.cache.GetAllTasks(userID)
	if err != nil && len(redisTasks) > 0 {
		return redisTasks, nil
	}

	//если в кеше нет ничего, беру из БД
	dbTasks, err := s.taskRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения задач из БД: %v", err)
	}

	return dbTasks, nil
}

func (s *TaskService) GetByID(taskID, userID uint) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID, userID)

	if err != nil {
		if errors.Is(err, app_errors.ErrRecordNotFound) {
			return nil, app_errors.ErrTaskNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return task, nil
}

func (s *TaskService) Update(task *models.Task) (int64, error) {
	// Проверяем существование задачи
	existingTask, err := s.taskRepo.GetByID(task.ID, task.UserID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при проверке задачи: %w", err)
	}

	if existingTask.DeletedAt.Valid {
		return 0, app_errors.ErrTaskDeleted
	}

	// Обновляем только разрешенные поля
	updateData := map[string]interface{}{
		"task_id":    task.ID,
		"user_id":    task.UserID,
		"completed":  task.Completed,
		"updated_at": time.Now(),
	}

	// if err != nil {
	// 	return 0, fmt.Errorf("ошибка сериализации данных: %w", err)
	// }

	// Добавляем в очередь обновлений
	queueSize, err := s.cache.PushUpdate(updateData)
	if err != nil {
		return 0, fmt.Errorf("ошибка сохранения в Redis: %w", err)
	}

	return queueSize, nil

}

func (s *TaskService) Delete(taskID, userID uint) (int64, error) {
	// Добавляем в очередь удалений
	deleteData := map[string]interface{}{
		"task_id": taskID,
		"user_id": userID,
	}
	queueSize, err := s.cache.PushDelete(deleteData)
	if err != nil {
		return 0, fmt.Errorf("ошибка сохранения в Redis: %w", err)
	}

	return queueSize, nil
}
