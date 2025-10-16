package service

import (
	"errors"
	"fmt"
	"pg-todolist/internal/app_errors"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"pg-todolist/pkg/utils"
	"pg-todolist/internal/dto"

)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*models.User, error) {
	//validate that user Does Not Exist
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, app_errors.ErrEmailExists
	}
	// password hashing
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка хеширования пароля: %w", err)
	}
	user:= &models.User{
		Email: req.Email,
		Password: hashedPassword,
	}

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
		return nil, fmt.Errorf("ошибка базы данных при поиске пользователя: %w", err)
	}
	//check password
	if !utils.CheckPassword(password, user.Password) {
		return nil, app_errors.ErrWrongPassword
	}

	return user, nil
}
