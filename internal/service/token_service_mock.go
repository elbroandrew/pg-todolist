package service

import "github.com/stretchr/testify/mock"



type TokenServiceMock struct {
	mock.Mock
}

func (m *TokenServiceMock) GenerateTokenPair(userID uint) (accessToken, refreshToken string, err error) {
	args := m.Called(userID)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *TokenServiceMock) ValidateAccessToken(tokenString string) (uint, error) {
	
	args := m.Called(tokenString)
	return uint(args.Int(0)), args.Error(1)
}

func (m *TokenServiceMock) RefreshTokens(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
	args := m.Called(oldRefreshToken)
	return args.String(0), args.String(1), args.Error(2)

}

func (m *TokenServiceMock) RevokeTokens(accessToken, refreshToken string) {
	m.Called(accessToken, refreshToken)
}
