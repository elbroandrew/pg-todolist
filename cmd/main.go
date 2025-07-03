package main

import (
	"context"
	"log"
	"os"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/router"
	"pg-todolist/internal/service"
	"pg-todolist/pkg/database"
	"strconv"

	"github.com/joho/godotenv"
)

func init(){

	if err := godotenv.Load(); err != nil {  
		log.Fatal("ERROR LOAD .ENV FILE")
	}
}

func main() {

	cfg := database.NewConfigFromEnv()
	db := database.InitMySQL(cfg)
	db.Connect(context.Background())

	defer db.Close()


	rdb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	//init REPOS
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	redisCacheRepo, _ := repository.NewRedisRepository(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		rdb,
	)

	defer redisCacheRepo.Close()

	//init SERVICES
	authService := service.NewAuthService(userRepo, redisCacheRepo)
	taskService := service.NewTaskService(taskRepo, redisCacheRepo)

	//init HANDLERS
	authHandler := handlers.NewAuthHandler(authService, redisCacheRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	//setup router
	r := router.SetupRouter(authHandler, taskHandler, redisCacheRepo)

	r.Run(":8080")

}
