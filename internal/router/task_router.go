package router

import (
	"net/http"
	"pg-todolist/internal/handlers"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Middleware для извлечения userID из заголовка, который устанавливает Gateway
func UserIdMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		userIDHeader := c.GetHeader("X-User-ID")
		if userIDHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "User ID not provided by gateway"})
			return 
		}
		// валидация ID
		userID, err := strconv.ParseUint(userIDHeader, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error":"Invalid user id format"})
			return 
		}

		c.Set("userID", uint(userID))
		c.Next()
	}
}

func SetupTaskServiceRouter(taskHandler *handlers.TaskHandler) *gin.Engine {
	//использую gin.New - чтоб наглядно видеть recovery
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(UserIdMiddleware())
	{
		tasksGroup.GET("", taskHandler.GetTasks)
		tasksGroup.POST("", taskHandler.CreateTask)
		tasksGroup.GET("/:id", taskHandler.GetTaskByID)
		tasksGroup.PUT("/:id", taskHandler.UpdateTask)
		tasksGroup.DELETE("/:id", taskHandler.DeleteTask)
	}

	//curl -X GET http://localhost:8081/health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Task Service is healthy"})
	})

	return r
}