package service

import (
	"errors"
	"fmt"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/interfaces"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"pg-todolist/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	userRepo *repository.UserRepository
	tokenRepo interfaces.TokenRepository
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo interfaces.TokenRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *AuthService) Register(user *models.User) (*models.User, error) {
	//validate that user Does Not Exist
	if _, err := s.userRepo.FindByEmail(user.Email); err == nil {
		return nil, app_errors.ErrEmailExists
	}
	// password hashing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}
	user.Password = hashedPassword

	//save user to DB
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return  user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, string, error) {
	// find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, app_errors.ErrRecordNotFound) {
			return nil, "", "", app_errors.ErrUserNotFound
		}
		return nil, "", "", fmt.Errorf("ошибка базы данных: %w", err)
	}
	
	//check password
	if !utils.CheckPassword(password, user.Password) {
		return nil, "", "", app_errors.ErrWrongPassword
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		return nil,"", "", fmt.Errorf("failed to generate tokens: %w", err)
	}
	// Сохраняем время последнего обновления
	if err := s.tokenRepo.SetLastRefresh(user.ID, time.Now()); err != nil {
		return nil,"", "", fmt.Errorf("failed to set last refresh: %w", err)
	}
	return user, accessToken, refreshToken, nil
}

func (s *AuthService) RevokeToken(token string) error {

	// Проверяем валидность токена перед отзывом
	if _, err := utils.ParseJWT(token); err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Добавляем в блеклист на 7 дней
	err := s.tokenRepo.RevokeToken(token, 24*7*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

func (s *AuthService) Logout(c *gin.Context) error {
    // Получаем токены
    accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
    refreshToken, _ := c.Cookie("refresh_token")

    // Отзыв access токена, если не истек
    if accessToken != "" {
        if claims, err := utils.GetTokenClaims(accessToken); err == nil {
            if exp, ok := claims["exp"].(float64); ok {
				expTime := time.Unix(int64(exp), 0)
				ttl := time.Until(expTime)
                if ttl > 0 {
                    go s.tokenRepo.RevokeToken(accessToken, ttl)
                }
			}
        }
    }

	// Отзыв refresh токена, если не истек
    if refreshToken != "" {
        if claims, err := utils.GetTokenClaims(refreshToken); err == nil {
            if exp, ok := claims["exp"].(float64); ok {
				expTime := time.Unix(int64(exp), 0)
				ttl := time.Until(expTime)
                if ttl > 0 {
                    go s.tokenRepo.RevokeToken(refreshToken, ttl)
                }
			}
        }
    }

    // Очищаем refresh token куки
    c.SetCookie("refresh_token", "", -1, "/", "", true, true)
    return nil
}