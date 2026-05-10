package github

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
)

type UserRepo interface {
	UpsertByGitHubID(ctx context.Context, ghUser dto.GitHubUser) (string, error)
}

type SessionStore interface {
	Create(ctx context.Context, userID string) (accessToken, refreshToken string, err error)
}
