package main

import (
	"log"
	"os"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/service"
	"pg-todolist/internal/router"
	"pg-todolist/pkg/database"
	"pg-todolist/pkg/server"

	// "github.com/joho/godotenv"
)



func main(){
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("ERROR LOAD .env FILE")
	// }

	db := database.InitMySQL()
	sqlDB, _ := db.DB()
	db.AutoMigrate(&models.Task{})  // task_service знает только о задачах.

	taskrepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskrepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	// router setup
	r := router.SetupTaskServiceRouter(taskHandler)

	port := os.Getenv("TASK_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	app := server.NewApp(r, port)
	// регистрирую методы для закрытия ресурсов
	app.AddCloser(func() error {
		log.Println("Closing TaskService MySQL connection...")
		return sqlDB.Close()
	})

	app.Run()
}