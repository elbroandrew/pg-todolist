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

		// Поулчаю access token
		accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		// если access токен валиден и не отозван - пропускаю запрос
		if accessToken != "" {
			if claims, err := utils.ParseJWTWithClaims(accessToken); err == nil {
				// если отозван
				if time.Until(time.Unix(int64(claims["exp"].(float64)), 0)) > 1*time.Minute {
					c.Set("userID", uint(claims["userID"].(float64)))
					c.Next()
					return
				}
				// проверяю блеклист только для истекшего/почти истекшего токена
				revoked, err := cache.IsTokenRevoked(accessToken)
				if err == nil && !revoked {
					c.Set("userID", uint(claims["userID"].(float64)))
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

		// Проверяем, не отозван ли refresh token (из блеклиста)
		revoked, err := cache.IsTokenRevoked(refreshToken)
		if err != nil || revoked {
			c.SetCookie("refresh_token", "", -1, "/", "", true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Сессия была отозвана",
				"code":  "session_revoked",
			})
			return
		}

		// Валидирую refresh токен
		claims, err := utils.ParseJWTWithClaims(refreshToken)
		if err != nil {
			// Автоматический отзыв истекшего refresh token
            if errors.Is(err, jwt.ErrTokenExpired) {
                go cache.RevokeToken(refreshToken, 24*time.Hour)
            }
			c.SetCookie("refresh_token", "", -1, "/", "", true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Сессия истекла.",
				"code":  "session_expired",
			})
			return
		}

		// Генерируем новую пару токен
		userID := uint(claims["userID"].(float64))
		
		newAccess, newRefresh, err := utils.GenerateTokens(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка обновления токенов.",
				"code":  "token_generation_failed",
			})
			return
		}
		//Проверка, когда последний раз обновлялся refresh token в Redis
		shouldRefresh := false
		lastRefresh, _ := cache.GetLastRefresh(userID)
		// Если было больше 24 часов назад обновление
		if time.Since(lastRefresh) > 24*time.Hour{
			shouldRefresh = true
		} else {
			// Ротация, если осталось меньше половины срока действия
			expTime := time.Unix(int64(claims["exp"].(float64)), 0) 
			if time.Until(expTime) < 24*time.Hour*7/2 {
				shouldRefresh = true
			} 
		}

		if shouldRefresh {
			// Отзываем старый refresh token перед установкой нового
            expTime := time.Unix(int64(claims["exp"].(float64)), 0)
            ttl := time.Until(expTime)
            if ttl > 0 {
                go cache.RevokeToken(refreshToken, ttl)
            }
            
            c.SetCookie("refresh_token", newRefresh, 3600*24*7, "/", "", true, true)
            cache.SetLastRefresh(userID, time.Now())
		// 	//Установка нового refresh токена
		// 	c.SetCookie(
		// 		"refresh_token",
		// 		newRefresh,
		// 		3600*24*7,
		// 		"/",
		// 		"",
		// 		true, // Secure
		// 		true) // HttpOnly
		// 	cache.SetLastRefresh(userID, time.Now())
		// } else {
		// 	// Используется существующий refresh token из куки
		// 	existingRefreshToken, _ := c.Cookie("refresh_token")
		// 	c.SetCookie(
		// 		"refresh_token",
		// 		existingRefreshToken,
		// 		3600*24*7,
		// 		"/",
		// 		"",
		// 		true,
		// 		true)
		 }

		c.Header("New-Access-Token", newAccess)
		c.Header("Access-Control-Expose-Headers", "New-Access-Token") // для CORS
		c.Set("userID", userID)

		c.Next()
	}

}
