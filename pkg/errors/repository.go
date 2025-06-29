package errors



var (
	ErrRecordNotFound 	= New(404, "record not found")
	ErrUserNotFound 	= New(404, "user not found")
	ErrNoRowsAffected 	= New(401, "записи в БД не были обновлены")
)