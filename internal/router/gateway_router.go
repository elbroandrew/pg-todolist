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
	"pg-todolist/internal/contextkeys"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)


//обработчик Reverse Proxy для проксирования запросов /tasks на TaskService. Он знает URL TaskService'a.
func NewReverseProxy(target string) gin.HandlerFunc {
	targetUrl, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request){

		//сначала вызову оригинальный директор, он правильно настроит req.URL.Host,  req.URL.Scheme,
		//и некоторые заголовки типа "X-Forwarded-For"
		originalDirector(req)


		userIDVal := req.Context().Value(contextkeys.UserIDKey)
		if userIDVal == nil {
			log.Println("userID not found in context for proxying")
			return 
		}
		userID, ok := userIDVal.(uint)
		if !ok {
			log.Println("userID in context is not of type uint")
			return 
		}
		//устанавливаю заголовок для внутреннего сервиса
		req.Header.Set("X-User-ID", strconv.FormatUint(uint64(userID), 10))
		req.Header.Del("Authorization")
		//явно указываю заголовок Host, чтобы целевой сервис правильно его видел, если он использует виртуальные хосты
		req.Host = targetUrl.Host
	}
	return func(c *gin.Context){
		if userID, exists := c.Get("userID"); exists {
			ctx := context.WithValue(c.Request.Context(), contextkeys.UserIDKey, userID)
			c.Request = c.Request.WithContext(ctx)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
	
}

func SetupGatewayRouter(
	authHandler *handlers.AuthHandler,
	tokenService *service.TokenService,
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
		c.JSON(http.StatusOK, gin.H{"status": "Gateway is ok"})
	})

	return r
}