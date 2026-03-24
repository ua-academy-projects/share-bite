package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *TokenManager) GeneratePair(userID string) (string, string, error) {
	accessToken, err := m.generate(userID, m.accessSecret, m.accessTTL)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := m.generate(userID, m.refreshSecret, m.refreshTTL)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (m *TokenManager) ParseRefreshToken(tokenStr string) (string, error) {
	return m.parse(tokenStr, m.refreshSecret)
}

func (m *TokenManager) ParseAccessToken(tokenStr string) (string, error) {
	return m.parse(tokenStr, m.accessSecret)
}

func (m *TokenManager) generate(userID, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (m *TokenManager) parse(tokenStr, secret string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid subject")
	}
	return sub, nil
}
