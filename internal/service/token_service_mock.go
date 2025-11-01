package service

import "github.com/stretchr/testify/mock"



type TokenServiceMock struct {
	mock.Mock
	ValidTokens map[string]uint
}

func NewTokenServiceMock() *TokenServiceMock {
    return &TokenServiceMock{
        ValidTokens: make(map[string]uint),
    }
}

func (m *TokenServiceMock) GenerateTokenPair(userID uint) (string, string, error) {
	args := m.Called(userID)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *TokenServiceMock) ValidateAccessToken(tokenString string) (uint, error) {
	// Если токен есть в валидных - возвращаем его
    if userID, exists := m.ValidTokens[tokenString]; exists {
        return userID, nil
    }
	args := m.Called(tokenString)
	if userID, ok := args.Get(0).(uint); ok {
        return userID, args.Error(1)
    }
    return 0, args.Error(1)
}

func (m *TokenServiceMock) RefreshTokens(oldRefreshToken string) (string, string, error) {
	args := m.Called(oldRefreshToken)
	return args.String(0), args.String(1), args.Error(2)

}

func (m *TokenServiceMock) RevokeTokens(accessToken, refreshToken string) error {
	delete(m.ValidTokens, accessToken)
	delete(m.ValidTokens, refreshToken)
	args := m.Called(accessToken, refreshToken)
	return args.Error(0)
}
