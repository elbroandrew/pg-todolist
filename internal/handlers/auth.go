package handlers

import (
	"net/http"
	"pg-todolist/internal/models"
	"pg-todolist/internal/service"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат почты"})
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
	var creds struct {
		Email 		string	`json:"email"`
		Password	string	`json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное тело запроса"})
		return
	}

	token, err := h.authService.Login(creds.Email, creds.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "неверный пароль" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"email": creds.Email,
		},
	})
}