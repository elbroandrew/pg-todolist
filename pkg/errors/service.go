package errors


var (
	ErrTaskNotFound 	= New(404, "user not found")
	ErrEmailExists 		= New(409, "email уже существует")
	ErrWrongPassword 	= New(401, "неверный пароль")
	ErrTaskDeleted      = New(404, "задача не существует")
)