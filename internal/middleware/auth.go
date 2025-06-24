package middleware

import (
	"log"
	"pg-todolist/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error":"Токен отсутствует"})
			return 
		}
		
		userID, err := utils.ParseJWT(token)
		if err != nil {
			// Логируем ошибку для отладки
            log.Printf("JWT validation error: %v", err)
			
			c.AbortWithStatusJSON(401, gin.H{
				"error":"Неверный токен",
				"details": err.Error(),
			})
			return 
		}
		c.Set("userID", userID)
		c.Next()
	}
}