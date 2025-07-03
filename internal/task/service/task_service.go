package service

import (
	"log"
	"pg-todolist/internal/models"
	"pg-todolist/internal/task/repository"
	"pg-todolist/pkg/app_errors"
)

type TaskService struct {
	repo repository.TaskRepository
	cache repository.CacheRepository
}

func NewTaskService(repo repository.TaskRepository, cache repository.CacheRepository) *TaskService {
	return &TaskService{
		repo: repo,
		cache: cache,
	}
}

func (s *TaskService) CreateTask(userID uint, title, description string) *app_errors.AppError {
	task := &models.Task{
		UserID:      userID,
		Title:       title,
		Completed:   false,
	}

	if err := s.repo.Create(task); err != nil {
		return err
	}

	// Инвалидируем кеш для этого пользователя
	s.cache.InvalidateUserTasks(userID)
	return nil
}

func (s *TaskService) GetUserTasks(userID uint) ([]models.Task, *app_errors.AppError) {
	// Пробуем получить из кеша
	if tasks, err := s.cache.GetUserTasks(userID); err == nil {
		return tasks, nil
	}

	// Если нет в кеше, идем в БД
	tasks, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	if err := s.cache.SetUserTasks(userID, tasks); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		log.Printf("Failed to cache tasks: %v", err)
	}

	return tasks, nil
}

func (s *TaskService) CompleteTask(taskID, userID uint) *app_errors.AppError {
	// Проверяем существование задачи
	_, err := s.repo.GetByID(taskID, userID)
	if err != nil {
		return err
	}

	// Обновляем статус
	if err := s.repo.Update(taskID, map[string]interface{}{
		"completed": true,
	}); err != nil {
		return err
	}

	// Инвалидируем кеш
	s.cache.InvalidateUserTasks(userID)
	return nil
}
