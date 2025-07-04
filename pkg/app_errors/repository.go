package app_errors

var (
	ErrRecordNotFound = New(404, "record not found")
	ErrUserNotFound   = New(404, "user not found")
	ErrNoRowsAffected = New(401, "записи в БД не были обновлены")
	ErrInternalServer = New(500, "internal server error")
	ErrTokenNotFound  = New(404, "token not found")
)
