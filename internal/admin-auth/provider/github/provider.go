package ghAuth

import (
	"context"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
)

type UserRepo interface {
	UpsertByGitHubID(ctx context.Context, ghUser dto.GitHubUser) (*dto.User, error)
}

type SessionStore interface {
	Create(ctx context.Context, userID string) (token string, err error)
}
