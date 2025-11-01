package service

import (
	"errors"
	"fmt"
	"log"
	"pg-todolist/pkg/cache"
	"pg-todolist/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrSessionRevoked   = errors.New("session has been revoked")
	ErrTokenRefreshLock = errors.New("token refresh is temporarily locked")
)

type TokenService struct {
	jwtSecret []byte
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{jwtSecret: []byte(secret)}
}

// GenerateTokenPair создает новую пару access и refresh токенов
func (s *TokenService) GenerateTokenPair(userID uint) (accessToken, refreshToken string, err error) {
	return utils.GenerateTokens(userID, s.jwtSecret)
}

// ValidateAccessToken проверяет access токен. Возвращает userID или ошибку.
func (s *TokenService) ValidateAccessToken(tokenString string) (uint, error) {
	// 1. Парсим токен
	claims, err := utils.ParseJWTWithClaims(tokenString, s.jwtSecret)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrTokenExpired
		}
		log.Printf("[AUTH_ERROR] Token validation failed: %v", err)
		return 0, ErrTokenInvalid
	}

	// 2. Проверяем, не отозван ли он (в черном списке)
	revoked, err := cache.IsTokenRevoked(tokenString)
	if err != nil {
		// Логируем ошибку Redis, но не блокируем пользователя
		log.Printf("[REDIS_ERROR] Redis error on blacklist check: %v\n", err)
	}
	if revoked {
		return 0, ErrSessionRevoked
	}

	userID := uint(claims["userID"].(float64))
	return userID, nil
}

// RefreshTokens обновляет пару токенов, используя refresh токен
func (s *TokenService) RefreshTokens(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// 1. Проверяем, не отозван ли старый refresh токен
	revoked, err := cache.IsTokenRevoked(oldRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("redis error on blacklist check: %w", err)
	}
	if revoked {
		return "", "", ErrSessionRevoked
	}

	// 2. Валидируем старый refresh токен
	claims, err := utils.ParseJWTWithClaims(oldRefreshToken, s.jwtSecret)
	if err != nil {
		// Если он истек, добавляем его в черный список на всякий случай
		if errors.Is(err, jwt.ErrTokenExpired) {
			go cache.RevokeToken(oldRefreshToken, 24*time.Hour)
			return "", "", ErrTokenExpired
		}
		return "", "", ErrTokenInvalid
	}

	userID := uint(claims["userID"].(float64))

	// 3. (Опционально) Проверяем блокировку на частое обновление
	lastRefresh, _ := cache.GetLastRefresh(userID)
	if time.Since(lastRefresh) < 5*time.Second { // Защита от спама
		return "", "", ErrTokenRefreshLock
	}

	// 4. Генерируем новую пару токенов
	newAccessToken, newRefreshToken, err = s.GenerateTokenPair(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// 5. Ротация refresh токена: старый добавляем в черный список, новый становится актуальным
	// Это повышает безопасность. Если refresh токен украли, он будет действовать только один раз.
	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	ttl := time.Until(expTime)
	if ttl > 0 {
		go cache.RevokeToken(oldRefreshToken, ttl)
	}
	cache.SetLastRefresh(userID, time.Now())

	return newAccessToken, newRefreshToken, nil
}

// RevokeTokens добавляет оба токена в черный список при логауте
func (s *TokenService) RevokeTokens(accessToken, refreshToken string) error {
	if accessToken != "" {
		if claims, err := utils.GetTokenClaims(accessToken, s.jwtSecret); err == nil {
			if exp, ok := claims["exp"].(float64); ok {
				ttl := time.Until(time.Unix(int64(exp), 0))
				if ttl > 0 {
					go cache.RevokeToken(accessToken, ttl)
				}
			}
		}
	}
	if refreshToken != "" {
		if claims, err := utils.GetTokenClaims(refreshToken, s.jwtSecret); err == nil {
			if exp, ok := claims["exp"].(float64); ok {
				ttl := time.Until(time.Unix(int64(exp), 0))
				if ttl > 0 {
					go cache.RevokeToken(refreshToken, ttl)
				}
			}
		}
	}

	return nil
}
