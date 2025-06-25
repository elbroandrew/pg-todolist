package app_errors

import "errors"

var (
	ErrTaskNotFound 	= errors.New("user not found")
	ErrEmailExists 		= errors.New("email уже существует")
	ErrWrongPassword 	= errors.New("неверный пароль")
)