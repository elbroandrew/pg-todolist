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
	"pg-todolist/pkg/server"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("ERROR LOAD .env FILE")
	}

	db := database.InitMySQL()
	sqlDB, _ := db.DB()
	db.AutoMigrate(&models.User{})  // Gateway знает только про пользователей

	redisCache := cache.InitRedis(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
	)


	//init REPOS -- Data Access Layer
	userRepo := repository.NewUserRepository(db)
	

	//init SERVICES -- Business Logic Layer
	authService := service.NewAuthService(userRepo)
	
	tokenService := service.NewTokenService()

	//init HANDLERS -- Presentation Layer
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	

	//setup router
	taskServiceURL := os.Getenv("TASK_SERVICE_URL")
	if taskServiceURL == "" {
		taskServiceURL = "http://localhost:8081"
	}
	r := router.SetupGatewayRouter(authHandler, tokenService, redisCache, taskServiceURL)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	app := server.NewApp(r, port)
	// регистрирую методы для закрытия ресурсов
	app.AddCloser(func() error {
		log.Println("Closing MySQL database connection...")
		return sqlDB.Close()
	})
	app.AddCloser(func() error {
		log.Println("Closing Redis connection...")
		return redisCache.Close()
	})
	app.Run()

}
