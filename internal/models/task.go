package models

import (

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title     string `json:"title" gorm:"not null"`
	Completed bool   `json:"completed" gorm:"default:false"`
	UserID    uint   `json:"user_id" gorm:"not null"`
}
