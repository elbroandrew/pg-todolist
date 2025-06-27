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
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) Create(task *models.Task) error {
	return s.taskRepo.Create(task)
}

func (s *TaskService) GetTaskByUserID(userID uint) ([]models.Task, error) {
	return s.taskRepo.GetByUserID(userID)
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

func (s *TaskService) Update(task *models.Task) error {
	// Проверяем существование задачи
	existingTask, err := s.taskRepo.GetByID(task.ID, task.UserID)
	if err != nil {
		return fmt.Errorf("ошибка при проверке задачи: %w", err)
	}

	if existingTask.DeletedAt.Valid {
		return app_errors.ErrTaskDeleted
	}

	// Обновляем только разрешенные поля
	updates := map[string]interface{}{
		"completed":  task.Completed,
		"updated_at": time.Now(),
	}

	//Обновляю задачу
	return s.taskRepo.Update(task.ID, updates)

}

func (s *TaskService) Delete(taskID, userID uint) error {
	return s.taskRepo.Delete(taskID, userID)
}
