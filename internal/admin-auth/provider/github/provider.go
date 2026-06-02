package github

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
)

type UserRepo interface {
	UpsertByGitHubID(ctx context.Context, ghUser dto.GitHubUser) (*dto.User, error)
	FindByID(ctx context.Context, userID string) (*dto.UserWithRole, error)
}

type SessionStore interface {
	Create(ctx context.Context, userID, role string) (token string, err error)
}
