package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



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