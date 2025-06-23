package service

import (
	"errors"
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
		return "", errors.New("пользователь уже сущестует")
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