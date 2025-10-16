package router

import (
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/middleware"
	"pg-todolist/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	taskHandler *handlers.TaskHandler,
	tokenService *service.TokenService,
	redisClient *redis.Client,
) *gin.Engine {
	r := gin.Default()

	//Health check
	//for HEAD only via curl : curl -I http://localhost:8080/health
	r.HEAD("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	//for GET request
	//curl -X GET http://localhost:8080/health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	//Global middleware (logging, CORS)
	r.Use(middleware.CORS())
	// r.Use(middleware.Logger())

	rateLimitMiddleware := middleware.NewRateLimiter(redisClient)

	//Группа для aутентификации
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", rateLimitMiddleware, authHandler.Register)
		authGroup.POST("/login", rateLimitMiddleware, authHandler.Login)
		authGroup.POST("/refresh", rateLimitMiddleware, authHandler.Refresh)
		authGroup.POST("/logout", middleware.AuthMiddleware(tokenService), authHandler.Logout)
		
	}
	//Группа для задач - требуется аутентификация
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(middleware.AuthMiddleware(tokenService)) // для JWT
	{
		tasksGroup.GET("", taskHandler.GetTasks)
		tasksGroup.POST("", taskHandler.CreateTask)
		tasksGroup.GET("/:id", taskHandler.GetTaskByID)
		tasksGroup.PUT("/:id", taskHandler.UpdateTask)
		tasksGroup.DELETE("/:id", taskHandler.DeleteTask)
	}

	return r

}
