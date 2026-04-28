package ghAuth

import "context"

type UserRepo interface {
	UpsertByGitHubID(ctx context.Context, ghUser GitHubUser) (*User, error)
}

type SessionStore interface {
	Create(ctx context.Context, userID string) (token string, err error)
}
