package repository

import (
	"context"
	"errors"
	"pg-todolist/internal/models"
	"pg-todolist/pkg/app_errors"
	"pg-todolist/pkg/database"

	"gorm.io/gorm"
)

type userRepository struct {
	db database.Database
}

// конструктор
func NewUserRepository(db database.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) *app_errors.AppError {
	if err := r.db.GetDB().WithContext(ctx).Create(user).Error; err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

// func (r *UserRepository) CreateWithTransaction(ctx context.Context, user *models.User) error {
// 	return r.db.WithTransaction(ctx, func(tx *gorm.DB) error {
// 		if err := tx.Create(user).Error; err != nil {
// 			return err
// 		}
// 		return tx.Model(user).Update("version", user.Version+1).Error
// 	})
// }

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, *app_errors.AppError) {
	var user models.User
	err := r.db.GetDB().WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, app_errors.ErrUserNotFound
	}
	if err != nil {
		return nil, app_errors.ErrDBError
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, *app_errors.AppError) {
	var user models.User
	err := r.db.GetDB().WithContext(ctx).
		First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, app_errors.ErrUserNotFound
	}
	if err != nil {
		return nil, app_errors.ErrDBError
	}
	return &user, nil
}
