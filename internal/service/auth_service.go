package service

import (
	"errors"
	"fmt"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"pg-todolist/pkg/cache"
	"pg-todolist/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
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

func (s *AuthService) Login(email, password string) (*models.User, error) {
	// find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, app_errors.ErrRecordNotFound) {
			return nil, app_errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("ошибка базы данных: %w", err)
	}
	//check password
	if !utils.CheckPassword(password, user.Password) {
		return nil, app_errors.ErrWrongPassword
	}

	return user, nil
}

func (s *AuthService) RevokeToken(token string) error {

	// Проверяем валидность токена перед отзывом
	if _, err := utils.ParseJWT(token); err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Добавляем в блеклист на 7 дней
	err := cache.RevokeToken(token, 24*7*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

func (s *AuthService) Logout(c *gin.Context) error {
    // 1. Получаем токены
    accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
    refreshToken, _ := c.Cookie("refresh_token")

    // 2. Отзываем оба токена
    if accessToken != "" {
        if claims, err := utils.GetTokenClaims(accessToken); err == nil {
            exp := time.Unix(int64(claims["exp"].(float64)), 0)
            _ = cache.RevokeToken(accessToken, time.Until(exp))
        }
    }

    if refreshToken != "" {
        _ = cache.RevokeToken(refreshToken, 24*time.Hour)
    }

    // 3. Очищаем куки
    c.SetCookie("refresh_token", "", -1, "/", "", true, true)
    return nil
}