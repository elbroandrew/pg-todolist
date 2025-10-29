package tests

import (
	"context"
)

type MockRedis struct{}

func (m *MockRedis) IsTokenRevoked(token string) (bool, error) {
	// Всегда возвращаем false, чтобы токены считались валидными
	return false, nil
}

func (m *MockRedis) RevokeToken(token string, expiration int64) error {
	return nil
}

func (m *MockRedis) Close() error {
	return nil
}