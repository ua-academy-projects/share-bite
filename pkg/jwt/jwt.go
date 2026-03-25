package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenManager struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenManager(accessSecret string, refreshSecret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *TokenManager) GenerateToken(userID string, role string) (string, string, error) {
	accessToken, err := m.generate(userID, role, m.accessSecret, m.accessTTL)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := m.generate(userID, role, m.refreshSecret, m.refreshTTL)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (m *TokenManager) ParseAccessToken(token string) (string, string, error) {
	return m.parse(token, m.accessSecret)
}
func (m *TokenManager) ParseRefreshToken(token string) (string, string, error) {
	return m.parse(token, m.refreshSecret)
}

func (m *TokenManager) generate(userID, role, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (m *TokenManager) parse(tokenStr, secret string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid claims")
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid subject")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid role")
	}
	return sub, role, nil
}
