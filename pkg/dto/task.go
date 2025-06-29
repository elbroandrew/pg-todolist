package dto

import "time"

// Создание задачи
type CreateTaskRequest struct {
    UserID      string `json:"user_id"`  // ID владельца
    Title       string `json:"title" validate:"required,max=100"`
    Description string `json:"description" validate:"max=500"`
}

// Обновление задачи
type UpdateTaskRequest struct {
    Title       string `json:"title" validate:"max=100"`
    Description string `json:"description" validate:"max=500"`
}

// Ответ с данными задачи (используется в API)
type TaskResponse struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}
