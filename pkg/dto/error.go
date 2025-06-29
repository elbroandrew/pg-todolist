package dto

// Стандартный ответ при ошибке
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}