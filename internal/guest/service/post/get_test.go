package post

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func TestPostService_Get_Succeeds(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", CreatedAt: now}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	out, err := svc.Get(context.Background(), "42", "customer-1")
	require.NoError(t, err)
	assert.Equal(t, "42", out.ID)
	assert.Equal(t, "42", repo.lastGetID)
	assert.Equal(t, "customer-1", repo.lastGetViewerID)
	assert.Equal(t, now, out.CreatedAt)
}

func TestPostService_Get_PropagatesRepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo get failed")
	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{}, repoErr
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	_, err := svc.Get(context.Background(), "42", "customer-1")
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}

func TestPostService_Get_WithMentions(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{ID: postID}, nil
		},
		listMentionsFn: func(ctx context.Context, postIDs []string) (map[string][]entity.PostMention, error) {
			return map[string][]entity.PostMention{
				"42": {
					{CustomerID: "user-1"},
					{CustomerID: "user-2"},
				},
			}, nil
		},
	}

	customerRepo := &customerRepoMock{
		getByIDsFn: func(ctx context.Context, ids []string) ([]entity.Customer, error) {
			return []entity.Customer{
				{ID: "user-1"},
				{ID: "user-2"},
			}, nil
		},
	}

	svc := New(
		repo,
		&venueProviderMock{},
		&followRepoMock{},
		customerRepo,
		&storageMock{},
		&txManagerMock{},
	)

	out, err := svc.Get(context.Background(), "42", "viewer")
	require.NoError(t, err)

	require.Len(t, out.Mentions, 2)
	assert.Equal(t, "user-1", out.Mentions[0].ID)
	assert.Equal(t, "user-2", out.Mentions[1].ID)
}

func TestPostService_Get_NoMentions(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{ID: postID}, nil
		},
		listMentionsFn: func(ctx context.Context, postIDs []string) (map[string][]entity.PostMention, error) {
			return map[string][]entity.PostMention{}, nil
		},
	}

	svc := New(
		repo,
		&venueProviderMock{},
		&followRepoMock{},
		&customerRepoMock{},
		&storageMock{},
		&txManagerMock{},
	)

	out, err := svc.Get(context.Background(), "42", "viewer")
	require.NoError(t, err)
	assert.Empty(t, out.Mentions)
}
