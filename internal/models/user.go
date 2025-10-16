package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"column:password_hash;not null"`
	Tasks    []Task `gorm:"foreignKey:UserID"`
}
