package handlers

import (
	"errors"
	"log"
	"net/http"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/interfaces"
	"pg-todolist/internal/models"

	"pg-todolist/internal/service"
	"pg-todolist/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	tokenCache interfaces.TokenRepository
}

func NewAuthHandler(authService *service.AuthService, tokenCache interfaces.TokenRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService, 
		tokenCache: tokenCache,
	}
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
    if err == nil && refreshToken != "" {
        claims, err := utils.ParseJWTWithClaims(refreshToken)
        if err == nil {
            // Проверяем срок действия
            if exp, ok := claims["exp"].(float64); ok {
                expTime := time.Unix(int64(exp), 0)
                if time.Now().Before(expTime) {
                    // Проверяем в Redis, не отозван ли токен
                    revoked, err := h.tokenCache.IsTokenRevoked(refreshToken)
                    if err != nil {
                        log.Printf("Cache error checking revoked token: %v", err)
                        // Продолжаем, считаем что токен не отозван
                    } else if !revoked {
                        c.JSON(http.StatusBadRequest, gin.H{
                            "error": "User already has active session",
                            "code":  "already_logged_in",
                        })
                        return
                    }
                }
            }
        }
    }
				
	var creds struct {
		Email 		string	`json:"email"`
		Password	string	`json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверное тело запроса",
		})
		return
	}

	user, accessToken, refreshToken, err  := h.authService.Login(creds.Email, creds.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, app_errors.ErrUserNotFound) ||
		   errors.Is(err, app_errors.ErrWrongPassword) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": "invalid credentials"})
		return
	}
	
	// Set refresh token in HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", true, true)

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
	
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}