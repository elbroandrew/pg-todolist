package router

import (
	"net/http"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/middleware"

	"github.com/gin-gonic/gin"
)



func SetupTaskServiceRouter(taskHandler *handlers.TaskHandler) *gin.Engine {
	//использую gin.New - чтоб наглядно видеть recovery
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// если пришел запрос на /tasks/, он не бyдет перенаправлен на /tasks
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(middleware.UserIdMiddleware())
	{
		tasksGroup.GET("", taskHandler.GetTasks)
		tasksGroup.POST("", taskHandler.CreateTask)
		tasksGroup.GET("/:id", taskHandler.GetTaskByID)
		tasksGroup.PUT("/:id", taskHandler.UpdateTask)
		tasksGroup.DELETE("/:id", taskHandler.DeleteTask)
	}

	//healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "TaskService is ok"})
	})

	return r
}