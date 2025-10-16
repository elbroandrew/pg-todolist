package main

import (
	"log"
	"os"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/models"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/router"
	"pg-todolist/internal/service"
	"pg-todolist/pkg/cache"
	"pg-todolist/pkg/database"

	"github.com/joho/godotenv"

)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("ERROR LOAD .env FILE")
	}

	db := database.InitMySQL()
	db.AutoMigrate(&models.User{}, &models.Task{})

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	redisClient := cache.InitRedis(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)

	defer cache.Close()

	//init REPOS -- Data Access Layer
	userRepo := repository.NewUserRepository(db)
	taskrepo := repository.NewTaskRepository(db)

	//init SERVICES -- Business Logic Layer
	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskrepo)
	tokenService := service.NewTokenService()

	//init HANDLERS -- Presentation Layer
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	taskHandler := handlers.NewTaskHandler(taskService)

	//setup router
	r := router.SetupRouter(authHandler, taskHandler, tokenService, redisClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("SERVER IS RUNNING ON PORT %s", port)
	if err := r.Run(":"+port); err != nil {
		log.Fatalf("FAILED TO RUN SERVER: %v", err)
	}

}
