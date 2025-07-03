package repository

import (
	"context"
	"pg-todolist/internal/models"
	"pg-todolist/pkg/app_errors"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) *app_errors.AppError
	FindByEmail(ctx context.Context, email string) (*models.User, *app_errors.AppError)
	FindByID(ctx context.Context, id uint) (*models.User, *app_errors.AppError)
	// Update(ctx context.Context, user *models.User) *app_errors.AppError
}

type TokenRepository interface {
	StoreRefreshToken(userID uint, token string, expiresAt time.Time) *app_errors.AppError
	GetRefreshToken(userID uint) (string, *app_errors.AppError)
	DeleteRefreshToken(userID uint) *app_errors.AppError
	AddToBlacklist(token string, expiresAt time.Time) *app_errors.AppError
	IsTokenBlacklisted(token string) (bool, *app_errors.AppError)
}