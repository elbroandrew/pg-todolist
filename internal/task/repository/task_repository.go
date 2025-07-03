package repository

import (
	"errors"
	"pg-todolist/internal/models"
	"pg-todolist/pkg/app_errors"
	"pg-todolist/pkg/database"

	"gorm.io/gorm"
)


type taskRepository struct {
	db database.Database
}

func NewTaskRepository(db database.Database) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) *app_errors.AppError {
	if err := r.db.GetDB().Create(task).Error; err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

func (r *taskRepository) GetByUserID(userID uint) ([]models.Task, *app_errors.AppError) {
	var tasks []models.Task
	err := r.db.GetDB().Where("user_id = ?", userID).Find(&tasks).Error
	if err != nil {
		 return nil, app_errors.ErrTaskNotFound
	}
	return tasks, nil
}

func (r *taskRepository) GetByID(taskID, userID uint) (*models.Task, *app_errors.AppError) {
	var task models.Task
	err := r.db.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, app_errors.ErrRecordNotFound
	}
	if err != nil {
		return nil, app_errors.ErrDBError
	}
	return &task, nil
}

func (r *taskRepository) Update(taskID uint, updates map[string]interface{}) *app_errors.AppError {
    result := r.db.GetDB().Model(&models.Task{}).
	Where("id = ? AND deleted_at IS NULL", taskID).
	Updates(updates)
	
	if result.Error != nil {
		return app_errors.ErrDBError
	}
	//Если была обновлена хотя бы одна запись
	if result.RowsAffected == 0 {
		return app_errors.ErrNoRowsAffected
	}
	return nil
}

func (r *taskRepository) Delete(taskID, userID uint) *app_errors.AppError {
	result := r.db.GetDB().
	Where("id = ? AND user_id = ?", taskID, userID).
	Delete(&models.Task{})

	if result.Error != nil {
		return app_errors.ErrDBError
	}
	if result.RowsAffected == 0 {
		return  app_errors.ErrTaskDeleted
	}

	return nil
}
