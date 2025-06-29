package main

import (
	"os"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/router"
	"pg-todolist/internal/service"
	"pg-todolist/pkg/database"
	"strconv"
)

func main() {

	db := database.InitMySQL()

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

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
