package repository

import (
	"errors"
	"fmt"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"

	"gorm.io/gorm"
)


type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetByUserID(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) GetByID(taskID, userID uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, app_errors.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &task, nil
}

func (r *TaskRepository) Update(taskID uint, updates map[string]interface{}) error {
    result := r.db.Model(&models.Task{}).
	Where("id = ? AND deleted_at IS NULL", taskID).
	Updates(updates)
	
	if result.Error != nil {
		return fmt.Errorf("ошибка БД при обновлении задачи: %w", result.Error)
	}
	//Если была обновлена хотя бы одна запись
	if result.RowsAffected == 0 {
		return app_errors.ErrNoRowsAffected
	}
	return nil
}

func (r *TaskRepository) Delete(taskID, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&models.Task{}).Error

}
