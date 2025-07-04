package main

import (
	"context"
	"log"
	userRepo "pg-todolist/internal/auth/repository"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/router"
	"pg-todolist/internal/auth/service"
	taskRepo "pg-todolist/internal/task/repository"
	cacheRepo "pg-todolist/internal/task/repository"
	tokenRepo "pg-todolist/internal/auth/repository"
	"pg-todolist/pkg/config"
	"pg-todolist/pkg/database"

	"github.com/joho/godotenv"
)

func init(){

	if err := godotenv.Load(); err != nil {  
		log.Fatal("ERROR LOAD .ENV FILE")
	}
}

func main() {

	cfg := config.NewConfigFromEnv()
	db := database.InitMySQL(cfg)
	db.Connect(context.Background())

	defer db.Close()

	redisCfg := config.NewRedisConfigFromEnv()
	redisClient, err := database.NewRedisClient(redisCfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisClient.Close()

	//init REPOS
	userRepo := userRepo.NewUserRepository(db)
	taskRepo := taskRepo.NewTaskRepository(db)
	cacheRepo := cacheRepo.NewTaskCache(redisClient)
	tokenRepo := tokenRepo.NewTokenRepository(redisClient)

	//init SERVICES
	authService := service.NewAuthService(userRepo, cacheRepo)
	taskService := service.NewTaskService(taskRepo, cacheRepo)

	//init HANDLERS
	authHandler := handlers.NewAuthHandler(authService, cacheRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	//setup router
	r := router.SetupRouter(authHandler, taskHandler, cacheRepo)

	r.Run(":8080")

}
