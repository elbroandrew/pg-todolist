package router

import (
	"log"
	"net/http"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/middleware"
	"pg-todolist/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//обработчик Reverse Proxy для проксирования запросов /tasks на TaskService. Он знает URL TaskService'a.
func NewReverseProxy(target string) gin.HandlerFunc {
	proxy, err := NewCustomReverseProxy(target)
    if err != nil {
        log.Fatalf("Failed to create reverse proxy: %v", err)
    }
    return proxy.GinHandler()
	
}

func SetupGatewayRouter(
	authHandler *handlers.AuthHandler,
	tokenService service.ITokenService,
	redisClient *redis.Client,
	taskServiceURL string,
) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.Use(middleware.CORS())


	rateLimitMiddleware := middleware.NewRateLimiter(redisClient)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", rateLimitMiddleware, authHandler.Register)
		authGroup.POST("/login", rateLimitMiddleware, authHandler.Login)
		authGroup.POST("/refresh", rateLimitMiddleware, authHandler.Refresh)
		authGroup.POST("/logout", middleware.AuthMiddleware(tokenService), authHandler.Logout)
		
	}

	//прокси эндпоинты
	tasksProxy := NewReverseProxy(taskServiceURL)
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(middleware.AuthMiddleware(tokenService))
	{
		// вот эти явные пути - чтоб Gin не делал редиректы на trailing slash в конце
		tasksGroup.GET("", tasksProxy)
		tasksGroup.POST("", tasksProxy)
		tasksGroup.DELETE("", tasksProxy)

		
		tasksGroup.Any("/*any", tasksProxy)
	}

	//healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Gateway is ok",
			"service": "gateway",
			"timestamp": time.Now().UTC(),
		})
	})

	return r
}