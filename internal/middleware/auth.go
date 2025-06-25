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
	return func(c *gin.Context) {

		// Пропускаем OPTIONS запросы
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		log.Printf("=> AuthMiddleware: %s", c.Request.URL.Path)
		// Try getting the access token
		accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		if accessToken != "" {
			if userID, err := utils.ParseJWT(accessToken); err == nil {
				c.Set("userID", userID)
				c.Next()
				return
			} else {
				// Логируем ошибку для отладки
				log.Printf("JWT validation error: %v", err)
			}

		}

		//если access token неавлиден - проверяю refresh token
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Требуется авторизация",
				"code":  "invalid_credentials",
			})
			return
		}

		// Валидируем refresh token
		userID, err := utils.ParseJWT(refreshToken)
		if err != nil {
			// Важно: инвалидируем cookie при невалидном refresh token
			c.SetCookie("refresh_token", "", -1, "/", "", true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Сессия истекла. Пожалуйста, войдите снова.",
				"code":  "session_expired",
			})
			return
		}

		// Проверяем, не отозван ли refresh token (блеклист)
		revoked, err := cache.IsTokenRevoked(refreshToken)
		if err != nil {
			log.Printf("Redis Error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка сервера.",
				"code":  "token_verification_failed",
			})
			return
		}

		if revoked {
			c.SetCookie("refresh_token", "", -1, "/", "", true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Сессия была отозвана",
				"code":  "session_revoked",
			})
			return
		}

		// Генерируем новую пару токенов
		newAccess, newRefresh, err := utils.GenerateTokens(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка обновления токенов.",
				"code":  "token_generation_failed",
			})
			return
		}

		// Обновляем токены
		c.SetCookie("refresh_token", newRefresh, 3600*24*7, "/", "", true, true) // Secure, HttpOnly
		c.Header("New-Access-Token", newAccess)
		c.Header("Access-Control-Expose-Headers", "New-Access-Token") // для CORS
		c.Set("userID", userID)

		// Добавляем старый refresh token в блеклист
		go func(oldToken string) {
			if err := cache.RevokeToken(oldToken, 24*7*time.Hour); err != nil {
				log.Printf("Failed to revoke token: %v", err)
			}
		}(refreshToken)

		log.Printf("<= AuthMiddleware (refresh token)")
	}

}
