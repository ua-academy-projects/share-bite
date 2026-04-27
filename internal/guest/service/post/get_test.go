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
	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

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
	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	_, err := svc.Get(context.Background(), "42", "customer-1")
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
