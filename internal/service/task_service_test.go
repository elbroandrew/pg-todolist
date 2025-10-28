package service

import (
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTaskService_GetByID_Success(t *testing.T) {
	mockRepo := new(repository.TaskRepositoryMock)
	taskService := NewTaskService(mockRepo)

	expectedTask := &models.Task{
		Model:  gorm.Model{ID: 1},
		Title:  "Test task",
		UserID: 10,
	}

	mockRepo.On("GetByID", uint(1), uint(10)).Return(expectedTask, nil)

	actualTask, err := taskService.GetByID(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, actualTask)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(repository.TaskRepositoryMock)
	taskService := NewTaskService(mockRepo)
	mockRepo.On("GetByID", uint(1), uint(10)).Return(nil, app_errors.ErrRecordNotFound)
	_, err := taskService.GetByID(1, 10)
	assert.Error(t, err)
	assert.Equal(t, app_errors.ErrTaskNotFound, err)
	mockRepo.AssertExpectations(t)
}
