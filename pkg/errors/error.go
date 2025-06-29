package errors


// Базовая структура для кастомных ошибок
type AppError struct {
    Code    int    // HTTP-статус код или внутренний код
    Message string // Человекочитаемое сообщение
    // Err     error  // Вложенная ошибка (опционально)
}

// func (e *AppError) Error() string {
//     if e.Err != nil {
//         return fmt.Sprintf("%s: %v", e.Message, e.Err)
//     }
//     return e.Message
// }

// Конструктор для создания ошибок
func New(code int, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        // Err:     err,
    }
}
