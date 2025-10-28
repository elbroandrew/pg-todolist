package repository

import "pg-todolist/internal/models"

type ITaskRepository interface {
	GetByUserID(userID uint) ([]models.Task, error)
	GetByID(taskID, userID uint) (*models.Task, error)
	Create(task *models.Task) error
	Update(taskID uint, updates map[string]interface{}) error
	Delete(taskID, userID uint) error
}
