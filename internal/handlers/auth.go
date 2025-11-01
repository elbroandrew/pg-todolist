package handlers

import (
	"errors"
	"log"
	"net/http"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/dto"
	"pg-todolist/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
	tokenService service.ITokenService
}

func NewAuthHandler(authService *service.AuthService, tokenService service.ITokenService) *AuthHandler {
	return &AuthHandler{authService: authService, tokenService: tokenService}
}


func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное тело запроса"})
		return
	}
	
	user, err := h.authService.Register(&req)
	if err != nil {
		if errors.Is(err, app_errors.ErrEmailExists){

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to register user."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message":"User registered successfully", "user_id":user.ID})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное тело запроса"})
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, app_errors.ErrUserNotFound) ||
		   errors.Is(err, app_errors.ErrWrongPassword) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": "invalid credentials"})
		return
	}
	accessToken, refreshToken, err := h.tokenService.GenerateTokenPair(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}
	// Set refresh token in HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: accessToken,
		User: dto.UserResponse{
			ID: user.ID,
			Email: user.Email,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context){
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token found."})
		return
	}
	newAccess, newRefresh, err := h.tokenService.RefreshTokens(refreshToken)
	if err != nil {
		if errors.Is(err, service.ErrSessionRevoked) || errors.Is(err, service.ErrTokenInvalid) {
			c.SetCookie("refresh_token", "", -1, "/", "", false, true)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", newRefresh, 3600*24*7, "/", "", false, true)
	c.JSON(http.StatusOK, dto.TokenResponse{AccessToken: newAccess})
	
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	refreshToken, _ := c.Cookie("refresh_token")

	//Отзываю оба токена
	if err := h.tokenService.RevokeTokens(accessToken, refreshToken); err != nil {
		log.Printf("Failed to revoke tokens: %v", err)
	}
	

	// Clear the refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}