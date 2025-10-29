package main

import (
	"log"
	"pg-todolist/internal/models"
	"pg-todolist/pkg/database"
)


func main(){
	log.Println("Starting databse migration...")

	db := database.InitMySQL()

	log.Println("Running auto-migration...")

	// здесь надо перечислить все модели из всех сервисов, которые должны быть в БД
	// порядок важен, если есть зависимости (User должен быть перед Task)
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
	)

	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)

	}
	log.Println("Databse migration completed sucessfully!")
}