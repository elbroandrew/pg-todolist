package main

import (
	"os"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/repository"
	"pg-todolist/internal/router"
	"pg-todolist/internal/service"
	"pg-todolist/pkg/cache"
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

	cache.InitRedis(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		rdb,
	)

	defer cache.Close()

	//init REPOS
	userRepo := repository.NewUserRepository(db)
	taskrepo := repository.NewTaskRepository(db)

	//init SERVICES
	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskrepo)

	//init HANDLERS
	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)

	//setup router
	r := router.SetupRouter(authHandler, taskHandler)

	r.Run(":8080")

}
