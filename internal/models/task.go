package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Completed bool   `gorm:"default:false"`
	UserID    uint   `gorm:"not null"`
}
