package repository

import "pg-todolist/internal/models"

type ITaskRepository interface {
	GetByUserID(userID uint) ([]models.Task, error)
	GetByID(taskID, userID uint) (*models.Task, error)
	Create(task *models.Task) error
	Update(taskID uint, updates map[string]interface{}) error
	Delete(taskID, userID uint) error
}

type IUserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
}
