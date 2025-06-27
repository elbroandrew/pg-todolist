package app_errors


import "errors"

var (
	ErrRecordNotFound 	= errors.New("record not found")
	ErrUserNotFound 	= errors.New("user not found")
	ErrNoRowsAffected 	= errors.New("записи в БД не были обновлены")
)