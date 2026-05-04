package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusMuted     UserStatus = "muted"
	UserStatusSuspended UserStatus = "suspended"
)

type CustomClaims struct {
	UserID string     `json:"sub"`
	Role   string     `json:"role"`
	Status UserStatus `json:"status,omitempty"`
	Type   string     `json:"type"`
	jwt.RegisteredClaims
}

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

func (m *TokenManager) GenerateToken(userID string, role string, status UserStatus) (string, string, error) {
	accessToken, err := m.generate(userID, role, status, m.accessSecret, m.accessTTL, "access")
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}
	refreshToken, err := m.generate(userID, role, "", m.refreshSecret, m.refreshTTL, "refresh")
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	return accessToken, refreshToken, nil
}

func (m *TokenManager) ParseAccessToken(token string) (string, string, UserStatus, error) {
	return m.parse(token, m.accessSecret, "access")
}

func (m *TokenManager) ParseRefreshToken(token string) (string, string, error) {
	userID, role, _, err := m.parse(token, m.refreshSecret, "refresh")
	if err != nil {
		return "", "", err
	}

	return userID, role, nil
}

func (m *TokenManager) GetRefreshTTL() time.Duration {
	return m.refreshTTL
}

func (m *TokenManager) generate(userID, role string, status UserStatus, secret string, ttl time.Duration, tokenType string) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	if ttl <= 0 {
		return "", errors.New("invalid jwt ttl")
	}

	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		Status: status,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (m *TokenManager) parse(tokenStr, secret string, expectedType string) (string, string, UserStatus, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", "", errors.New("token is expired")
		}
		return "", "", "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", "", "", errors.New("invalid token")
	}
	if claims.Type != expectedType {
		return "", "", "", fmt.Errorf("invalid token type: %v", claims.Type)
	}

	return claims.UserID, claims.Role, claims.Status, nil
}
