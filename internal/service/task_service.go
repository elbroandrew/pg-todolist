package service

import (
	"errors"
	"fmt"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
)
var ErrTaskNotFound = errors.New("user not found")

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
		if errors.Is(err, repository.ErrRecordNotFound){
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return task, nil
}

func (s *TaskService) Update(task *models.Task) error {
	return s.taskRepo.Update(task)
}

func (s *TaskService) Delete(taskID, userID uint) error {
	return s.taskRepo.Delete(taskID, userID)
}
