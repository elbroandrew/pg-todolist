package router

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"pg-todolist/internal/handlers"
	"pg-todolist/internal/middleware"
	"pg-todolist/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//обработчик Reverse Proxy для проксирования запросов /tasks на TaskService. Она знает URL TaskService'a.
func NewReverseProxy(target string) gin.HandlerFunc {
	targetUrl, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	proxy.Director = func(req *http.Request){
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.URL.Path = req.URL.Path
		req.Host = targetUrl.Host

		//использую типизированный ключ для контекста, чтобы избежать коллизий
		type contextKey string
		const userIDKey contextKey = "userID"

		userIDVal := req.Context().Value(userIDKey)
		if userIDVal == nil {
			log.Println("userID not found in context for proxying")
			return 
		}
		userID, ok := userIDVal.(uint)
		if !ok {
			log.Println("userID in context is not of type uint")
			return 
		}
		req.Header.Set("X-User-ID", strconv.FormatUint(uint64(userID), 10))
		req.Header.Del("Authorization")
	}
	return func(c *gin.Context){
		type contextKey string
		const userIDKey contextKey = "userID"
		if userID, exists := c.Get("userID"); exists {
			ctx := context.WithValue(c.Request.Context(), userIDKey, userID)
			c.Request = c.Request.WithContext(ctx)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
	
}

func SetupGatewayRouter(
	authHandler *handlers.AuthHandler,
	tokenService *service.TokenService,
	redisClient *redis.Client,
	taskSrviceURL string,
) *gin.Engine {
	r := gin.Default()
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
	tasksProxy := NewReverseProxy(taskSrviceURL)
	tasksGroup := r.Group("/tasks")
	tasksGroup.Use(middleware.AuthMiddleware(tokenService))
	{
		tasksGroup.Any("/*any", tasksProxy)
	}

	//healthcheck
	//curl -X GET http://localhost:8080/health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Gateway is ok"})
	})

	return r
}