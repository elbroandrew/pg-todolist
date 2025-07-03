package repository

import (
	"context"
	"pg-todolist/internal/models"
	"pg-todolist/pkg/app_errors"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) *app_errors.AppError
	FindByEmail(ctx context.Context, email string) (*models.User, *app_errors.AppError)
	FindByID(ctx context.Context, id uint) (*models.User, *app_errors.AppError)
	// Update(ctx context.Context, user *models.User) *app_errors.AppError
}