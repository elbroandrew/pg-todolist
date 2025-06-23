package router

import (
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHadler *handlers.AuthHandler,
	taskHandler *handlers.TaskHandler,
) *gin.Engine {
	r := gin.Default()

	//Global middleware (logging, CORS)
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	//Группа для вутентификации
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHadler.Register)
		authGroup.POST("/login", authHadler.Login)
	}
	//Группа для задач - требуется аутентификация
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(middleware.AuthMiddleware()) // JWT
	{
		tasksGroup.GET("/", taskHandler.GetTasks)
		tasksGroup.POST("/", taskHandler.CreateTask)
		tasksGroup.GET("/:id", taskHandler.GetTaskByID)
		tasksGroup.PUT("/:id", taskHandler.UpdateTask)
		tasksGroup.DELETE("/:id", taskHandler.DeleteTask)
	}

	//Health check
	r.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{"status":"ok"})
	})

	return r

}
