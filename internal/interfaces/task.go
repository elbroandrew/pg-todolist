package interfaces

import "pg-todolist/internal/models"

type TaskService interface {
    Create(task *models.Task) error
    GetTaskByUserID(userID uint) ([]models.Task, error)
    GetByID(taskID, userID uint) (*models.Task, error)
    Update(task *models.Task) (int64, error) 
    Delete(taskID, userID uint) error
}