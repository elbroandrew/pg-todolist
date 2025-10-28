package repository

import (
	"pg-todolist/internal/models"

	"github.com/stretchr/testify/mock"
)


type TaskRepositoryMock struct {
	mock.Mock
}

func (m *TaskRepositoryMock) GetByUserID(userID uint) ([]models.Task, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *TaskRepositoryMock) GetByID(taskID, userID uint) (*models.Task, error) {
	args := m.Called(taskID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}


func (m *TaskRepositoryMock) Create(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0) 
}

func (m *TaskRepositoryMock) Update(taskID uint, updates map[string]interface{}) error {
	args := m.Called(taskID, updates)
	return args.Error(0)
}

func (m *TaskRepositoryMock) Delete(taskID, userID uint) error {
	args := m.Called(taskID, userID)
	return args.Error(0)
}