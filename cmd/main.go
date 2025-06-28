package main

import (
	"fmt"
	"os"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/repository/cache"
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

	rdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisCacheRepo, err := cache.NewRedisRepository(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		rdb,
	)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	defer redisCacheRepo.Close()

	//init REPOS
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	//init SERVICES
	authService := service.NewAuthService(userRepo, redisCacheRepo)
	taskService := service.NewTaskService(taskRepo, redisCacheRepo)

	//init HANDLERS
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)

	//setup router
	r := router.SetupRouter(authHandler, taskHandler)

	r.Run(":8080")

}
