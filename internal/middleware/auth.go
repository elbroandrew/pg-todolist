package middleware

import (
	"log"
	"net/http"
	"pg-todolist/pkg/cache"
	"pg-todolist/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		// Try getting the access token
		accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		
		if userID, err := utils.ParseJWT(accessToken); err == nil {
			// Логируем ошибку для отладки
            log.Printf("JWT validation error: %v", err)
		
			c.Set("userID", userID)
			c.Next() 
			return 
		}

		//если access token неавлиден - проверяю refresh token
		refreshToken, err := c.Cookie("refresh_token")
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Требуется авторизация",
                "code": "invalid_credentials",
            })
            return
        }
	
		// Проверяем, не отозван ли refresh token (добавьте эту проверку если используете блеклист)
        revoked, err := cache.IsTokenRevoked(refreshToken)
		if err != nil {
            
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "error": "Ошибка сервера.",
                "code": "token_verification_failed",
            })
            return
        }

		if revoked {
			c.SetCookie("refresh_token", "", -1, "/", "", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Сессия была отозвана",
				"code": "session_revoked",
			})
			return 
		}
        // Валидируем refresh token
        userID, err := utils.ParseJWT(refreshToken)
        if err != nil {
            // Важно: инвалидируем cookie при невалидном refresh token
            c.SetCookie("refresh_token", "", -1, "/", "", false, true)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Сессия истекла. Пожалуйста, войдите снова.",
                "code": "session_expired",
            })
            return
        }

		  // Генерируем новую пару токенов
        newAccess, newRefresh, err := utils.GenerateTokens(userID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "error": "Ошибка обновления токенов.",
                "code": "token_generation_failed",
            })
            return
        }

        // Обновляем токены
        c.SetCookie("refresh_token", newRefresh, 3600*24*7, "/", "", false, true) // Secure, HttpOnly
        c.Header("New-Access-Token", newAccess)
        c.Header("Access-Control-Expose-Headers", "New-Access-Token") // Важно для CORS
        c.Set("userID", userID)

        // Добавляем старый refresh token в блеклист 
        go cache.RevokeToken(refreshToken, 24*7*time.Hour) 

        c.Next()
    }
		
}