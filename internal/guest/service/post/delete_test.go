package post

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

func TestPostService_Delete_SucceedsForOwner(t *testing.T) {
	t.Parallel()

	deleted := false
	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusPublished, CreatedAt: time.Now()}, nil
		},
		updateStatusFn: func(ctx context.Context, postID, customerID string, status entity.PostStatus) error {
			deleted = true
			assert.Equal(t, "42", postID)
			assert.Equal(t, "customer-1", customerID)
			assert.Equal(t, entity.PostStatusDeleted, status)
			return nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	err := svc.Delete(context.Background(), "42", "customer-1")
	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestPostService_Delete_RejectsForeignOwner(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusPublished, CreatedAt: time.Now()}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	err := svc.Delete(context.Background(), "42", "customer-2")
	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code.NotFound, appErr.Code)
}

func TestPostService_Delete_IsIdempotentForAlreadyDeleted(t *testing.T) {
	t.Parallel()

	updateCalled := false
	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusDeleted, CreatedAt: time.Now()}, nil
		},
		updateStatusFn: func(ctx context.Context, postID, customerID string, status entity.PostStatus) error {
			updateCalled = true
			return nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	err := svc.Delete(context.Background(), "42", "customer-1")
	require.NoError(t, err)
	assert.False(t, updateCalled)
}
