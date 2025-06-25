package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret []byte

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not found in .env file")
	}

	jwtSecret = []byte(secret)
}

func GenerateTokens(userID uint) (accessToken, refreshToken string, err error){
	accessToken, err = GenerateJWT(userID, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateJWT(userID, 24*7*time.Hour)
	return accessToken, refreshToken, err
}

func GenerateJWT(userID uint, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(duration).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (uint, error) {
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
