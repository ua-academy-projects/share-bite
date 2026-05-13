package env

import (
	"fmt"
	"strconv"
)

type AuthConfig struct {
	maxSessions int
}

func NewAuthConfig(opts ...Options) (*AuthConfig, error) {
	val := GetSecret("AUTH_MAX_SESSIONS")
	if val == "" {
		return nil, fmt.Errorf("AUTH_MAX_SESSIONS is required but not set in environment")
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		return nil, fmt.Errorf("AUTH_MAX_SESSIONS must be an integer: %w", err)
	}

	if res <= 0 {
		return nil, fmt.Errorf("AUTH_MAX_SESSIONS must be a positive integer, got: %d", res)
	}

	return &AuthConfig{maxSessions: res}, nil
}

func (c *AuthConfig) MaxSessions() int {
	return c.maxSessions
}
