package repository

import (
	"gorm.io/gorm"
	"pg-todolist/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}
// ф-ия конструктор
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}