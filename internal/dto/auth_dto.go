package dto

// RegisterRequest DTO для регистрации пользователя
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
}

// LoginRequest DTO для входа пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse DTO для ответа с токенами
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// UserResponse DTO для информации о пользователе в ответе
type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// LoginResponse DTO для полного ответа при логине
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}