package dto

import "pg-todolist/internal/models"

// CreateTaskRequest DTO для создания задачи
type CreateTaskRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateTaskRequest DTO для обновления задачи
type UpdateTaskRequest struct {
	Completed *bool `json:"completed" binding:"required"`
}

// TaskResponse DTO для ответа с задачей
type TaskResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// FromModel конвертирует модель Task в TaskResponse DTO
func TaskResponseFromModel(task *models.Task) TaskResponse {
	return TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// TasksResponseFromModels конвертирует срез моделей в срез DTO
func TasksResponseFromModels(tasks []models.Task) []TaskResponse {
	responses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = TaskResponseFromModel(&task)
	}
	return responses
}