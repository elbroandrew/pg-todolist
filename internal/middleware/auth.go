package middleware

import (
	"net/http"
	"pg-todolist/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {

		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"error":"Authorization header is required."})
			return 
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format."})
			return
		}

		userID, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			if err == service.ErrTokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token has expired",
					"code": "token_expired",
				})
				return 
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":err.Error()})
			return 
		}

		c.Set("userID", userID)

		c.Next()
	}

}
