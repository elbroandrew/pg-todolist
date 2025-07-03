package main

import (
	"context"
	"log"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/router"
	"pg-todolist/internal/service"
	"pg-todolist/internal/task/repository"
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
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	cacheRepo := repository.NewTaskCache(redisClient)

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
