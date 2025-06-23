package service

import (
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
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

func (s *TaskService) Update(task *models.Task) error {
	return s.taskRepo.Update(task)
}

func (s *TaskService) Delete(taskID, userID uint) error {
	return s.taskRepo.Delete(taskID, userID)
}
