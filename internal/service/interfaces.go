package service

type ITokenService interface {
	GenerateTokenPair(userID uint) (accessToken, refreshToken string, err error) 
	ValidateAccessToken(tokenString string) (uint, error) 
	RefreshTokens(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error)
	RevokeTokens(accessToken, refreshToken string)
}