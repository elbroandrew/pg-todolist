package request

type Register struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}