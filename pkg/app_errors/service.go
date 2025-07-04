package app_errors


var (
	ErrTaskNotFound 	= New(404, "task not found")
	ErrEmailExists 		= New(409, "email уже существует")
	ErrWrongPassword 	= New(401, "неверный пароль")
	ErrTaskDeleted      = New(404, "задача не существует")
)