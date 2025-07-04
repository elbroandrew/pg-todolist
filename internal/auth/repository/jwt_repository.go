package repository

import (
	"context"
	"fmt"
	"pg-todolist/pkg/app_errors"
	"pg-todolist/pkg/database"
	"time"

	"github.com/redis/go-redis/v9"
)


type tokenRepository struct {
	client database.RedisClient
}

func NewTokenRepository(client database.RedisClient) TokenRepository {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) StoreRefreshToken(userID uint, token string, expiresAt time.Time) *app_errors.AppError {
	err := r.client.Set(context.Background(), 
		fmt.Sprintf("refresh_tokens:%d", userID), 
		token, 
		time.Until(expiresAt),
	)
	
	if err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

func (r *tokenRepository) GetRefreshToken(userID uint) (string, *app_errors.AppError) {
	token, err := r.client.Get(context.Background(), 
		fmt.Sprintf("refresh_tokens:%d", userID),
	)

	if err == redis.Nil {
		return "", app_errors.ErrTokenNotFound
	}
	if err != nil {
		return "", app_errors.ErrDBError
	}
	return token, nil
}

func (r *tokenRepository) DeleteRefreshToken(userID uint) *app_errors.AppError {
	err := r.client.Del(context.Background(), 
		fmt.Sprintf("refresh_tokens:%d", userID),
	)

	if err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

func (r *tokenRepository) AddToBlacklist(token string, expiresAt time.Time) *app_errors.AppError {
	err := r.client.Set(context.Background(), 
		fmt.Sprintf("blacklist:%s", token), 
		"1", 
		time.Until(expiresAt),
	)

	if err != nil {
		return app_errors.ErrDBError
	}
	return nil
}

func (r *tokenRepository) IsTokenBlacklisted(token string) (bool, *app_errors.AppError) {
	exists, err := r.client.Exists(context.Background(), 
		fmt.Sprintf("blacklist:%s", token),
	)

	if err != nil {
		return false, app_errors.ErrDBError
	}
	return exists == 1, nil
}