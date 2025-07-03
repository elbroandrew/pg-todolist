package app_errors


// Базовая структура для кастомных ошибок
type AppError struct {
	Code    int    // HTTP-статус код или внутренний код
	Message string // Человекочитаемое сообщение
}


// Конструктор для создания ошибок
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
