package repository

import (
	"pg-todolist/internal/models"
	"pg-todolist/pkg/app_errors"
)

type TaskRepository interface {
	Create(task *models.Task) *app_errors.AppError
	GetByUserID(userID uint) ([]models.Task, *app_errors.AppError)
	GetByID(taskID, userID uint) (*models.Task, *app_errors.AppError)
	Update(taskID uint, updates map[string]interface{}) *app_errors.AppError 
	Delete(taskID, userID uint) *app_errors.AppError
}

type CacheRepository interface {
	SetUserTasks(userID uint, tasks []models.Task) *app_errors.AppError
	GetUserTasks(userID uint) ([]models.Task, *app_errors.AppError)
	InvalidateUserTasks(userID uint) *app_errors.AppError
	MarkCompleted(taskID uint) *app_errors.AppError
	MarkDeleted(taskID uint) *app_errors.AppError
}