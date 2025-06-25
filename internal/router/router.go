package router

import (
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	taskHandler *handlers.TaskHandler,
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
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	//Группа для aутентификации
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)

		/*****TEST TOKEN*****/
		authGroup.GET("/validate", authHandler.ValidateToken)
		/****/
	}
	//Группа для задач - требуется аутентификация
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(middleware.AuthMiddleware()) // JWT
	{
		tasksGroup.GET("", taskHandler.GetTasks)
		tasksGroup.POST("", taskHandler.CreateTask)
		tasksGroup.GET("/:id", taskHandler.GetTaskByID)
		tasksGroup.PUT("/:id", taskHandler.UpdateTask)
		tasksGroup.DELETE("/:id", taskHandler.DeleteTask)
	}

	return r

}
