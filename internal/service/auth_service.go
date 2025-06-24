package service

import (
	"errors"
	"fmt"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"pg-todolist/pkg/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(user *models.User) (string, error) {
	//validate that user Does Not Exist
	_, err := s.userRepo.FindByEmail(user.Email)
	if err == nil {
		return "", errors.New("пользователь уже существует")
	}
	// password hashing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashedPassword

	//save user to DB
	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	//Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	return token, err
}

func (s *AuthService) Login(email, password string) (string, error) {
	// find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, app_errors.ErrUserNotFound) {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("database error: %w", err)
	}
	//check password
	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("неверный пароль")
	}
	//generate jwt token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", fmt.Errorf("ошибка создания токена: %w", err)
	}
	return token, nil
}
