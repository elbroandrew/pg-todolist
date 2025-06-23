package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (uint, error) {
	// 1. parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signature algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	// 2. check errors
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга токена: %v", err)
	}
	// 3. fetch the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//check expiration time
		exp, err := claims.GetExpirationTime()
		if err != nil || exp.Before(time.Now()) {
			return 0, fmt.Errorf("токен просрочен")
		}
		//getting the user ID
		userID, ok := claims["userID"].(float64)
		if !ok {
			return 0, fmt.Errorf("неверный формат userID в токене")
		}
		return uint(userID), nil
	}
	return 0, fmt.Errorf("недействительный токен")
}
