package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



func GenerateTokens(userID uint, jwtSecret []byte) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateJWT(userID, 15*time.Minute, jwtSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateJWT(userID, 24*7*time.Hour, jwtSecret)
	return accessToken, refreshToken, err
}

func GenerateJWT(userID uint, duration time.Duration, jwtSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(duration).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string, jwtSecret []byte) (uint, error) {
	// обрезаю "Bearer "
	if after, ok := strings.CutPrefix(tokenString, "Bearer "); ok {
		tokenString = after
	}
	// parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signature algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	// check errors
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга токена: %w", err)
	}
	// fetch the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check expiration time
		exp, err := claims.GetExpirationTime()
		if err != nil || exp.Before(time.Now()) {
			return 0, fmt.Errorf("токен просрочен")
		}
		// getting the user ID
		userID, ok := claims["userID"].(float64)
		if !ok {
			return 0, fmt.Errorf("неверный формат userID в токене")
		}
		return uint(userID), nil
	}
	return 0, fmt.Errorf("недействительный токен")
}

func GetTokenClaims(tokenString string, jwtSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func ParseJWTWithClaims(tokenString string, jwtSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token is malformed: %w", err) // неверный формат токена
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, fmt.Errorf("token signature is invalid: %w", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, err
		}
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
