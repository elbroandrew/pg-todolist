package handlers

import (
	"errors"
	"net/http"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/service"
	"pg-todolist/pkg/cache"
	"pg-todolist/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
        return
    }

    _, err := utils.ParseJWT(token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "valid": false,
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{"valid": true})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное тело запроса"})
		return
	}
	// validate email & password
	if !utils.ValidateEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный email"})
		return
	}
	if !utils.ValidatePassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль должен быть не менее 3 символов"})
		return
	}

	token, err := h.authService.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token":token})
}

func (h *AuthHandler) Login(c *gin.Context) {

	 // Проверяем, есть ли уже валидный refresh token
    refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		revoked, _ := cache.IsTokenRevoked(refreshToken)
		if !revoked{
			if _, err := utils.ParseJWTWithClaims(refreshToken); err == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Пользователь уже авторизован.",
					"code":  "already_logged_in",
				})
				return
			}
		}
    }

	var creds struct {
		Email 		string	`json:"email"`
		Password	string	`json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное тело запроса"})
		return
	}

	user, err := h.authService.Login(creds.Email, creds.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, app_errors.ErrUserNotFound) ||
		   errors.Is(err, app_errors.ErrWrongPassword) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": "invalid credentials"})
		return
	}
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}
	// Set refresh token in HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id": user.ID,
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.authService.Logout(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Logout failed",
			"code": "logout_failed",
		})
	 	return
	
	}
	// refreshToken, err := c.Cookie("refresh_token")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token"})
	// 	return
	// }

	// if err := h.authService.RevokeToken(refreshToken); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
	// 	return
	// }

	// // Clear the refresh token cookie
	// c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}