package middleware

import (
	"errors"
	"net/http"
	"pg-todolist/pkg/cache"
	"pg-todolist/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Try getting the access token
		accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		if accessToken != "" {
			if userID, err := utils.ParseJWT(accessToken); err == nil {
				// Check if the access token has been revoked
				revoked, err := cache.IsTokenRevoked(accessToken)
				if err == nil && !revoked {
					c.Set("userID", userID)
					c.Next()
					return
				}
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

		// Проверяем, не отозван ли refresh token (блеклист)
		revoked, err := cache.IsTokenRevoked(refreshToken)
		if err != nil {
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

		// Валидация refresh token
		// userID, err := utils.ParseJWT(refreshToken)
		// if err != nil {
		// 	// Важно: инвалидируем cookie при невалидном refresh token
		// 	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Сессия истекла. Пожалуйста, войдите снова.",
		// 		"code":  "session_expired",
		// 	})
		// 	return
		// }
		claims, err := utils.ParseJWTWithClaims(refreshToken)
        if err != nil {
            // Если refresh token истек - добавляем в blacklist
            if errors.Is(err, jwt.ErrTokenExpired) {
                go cache.RevokeToken(refreshToken, 24*time.Hour)
            }
            c.SetCookie("refresh_token", "", -1, "/", "", true, true)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Session expired",
                "code": "session_expired",
            })
            return
        }

		// Генерируем новую пару токенов
		userID := uint(claims["userID"].(float64))
		newAccess, newRefresh, err := utils.GenerateTokens(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка обновления токенов.",
				"code":  "token_generation_failed",
			})
			return
		}

		// Отзыв старых токенов
		// go func() {
		// 	// отзываю старый рефреш токен на 24 часа
		// 	_ = cache.RevokeToken(refreshToken, 24*time.Hour)
		// 	// отзываю старый access token на оставшееся время его жизни
		// 	if claims, err := utils.GetTokenClaims(accessToken); err == nil {
		// 		exp := time.Unix(int64(claims["exp"].(float64)), 0)
		// 		_ = cache.RevokeToken(accessToken, time.Until(exp))
		// 	}
		// }()

		// Обновляем токены
		c.SetCookie(
			"refresh_token",
			newRefresh,
			3600*24*7,
			"/",
			"",
			true, // Secure
			true) // HttpOnly
		c.Header("New-Access-Token", newAccess)
		c.Header("Access-Control-Expose-Headers", "New-Access-Token") // для CORS
		c.Set("userID", userID)

		c.Next()
	}

}
