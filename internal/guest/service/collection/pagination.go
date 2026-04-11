package collection

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func (s *service) generatePageToken(createdAt time.Time, id string) string {
	timeStr := createdAt.Format(time.RFC3339Nano)
	raw := fmt.Sprintf("%s|%s", timeStr, id)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func (s *service) parsePageToken(token string) (time.Time, string, error) {
	if token == "" {
		return time.Time{}, "", nil
	}

	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("invalid token encoding: %w", err)
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return time.Time{}, "", fmt.Errorf("invalid token format")
	}

	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, "", fmt.Errorf("invalid token time format: %w", err)
	}

	return createdAt, parts[1], nil
}
