package post

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

func TestPostService_Like_Succeeds(t *testing.T) {
	t.Parallel()

	liked := false
	published := make(chan bool, 1)
	now := time.Now()

	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "author-1", CreatedAt: now}, nil
		},
		likeFn: func(ctx context.Context, postID string, customerID string) error {
			liked = true
			return nil
		},
		getAuthorUserIDFn: func(ctx context.Context, postID string) (string, error) {
			return "user-1", nil
		},
	}

	publisher := &publisherMock{
		publishFn: func(ctx context.Context, target string, msg notification.Message) error {
			assert.Equal(t, "user-1", target)
			assert.Equal(t, notification.PostLiked, msg.Type)
			published <- true
			return nil
		},
	}

	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{}, WithPublisher(publisher))

	err := svc.Like(context.Background(), "42", "customer-1")
	require.NoError(t, err)
	assert.True(t, liked)

	select {
	case <-published:
		// success
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for publish")
	}
}

func TestPostService_Like_NoNotificationForSelfLike(t *testing.T) {
	t.Parallel()

	liked := false
	published := false
	fetched := false
	now := time.Now()

	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			fetched = true
			return entity.Post{ID: postID, CustomerID: "customer-1", CreatedAt: now}, nil
		},
		likeFn: func(ctx context.Context, postID string, customerID string) error {
			liked = true
			return nil
		},
		getAuthorUserIDFn: func(ctx context.Context, postID string) (string, error) {
			return "user-1", nil
		},
	}

	publisher := &publisherMock{
		publishFn: func(ctx context.Context, target string, msg notification.Message) error {
			published = true
			return nil
		},
	}

	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{}, WithPublisher(publisher))

	err := svc.Like(context.Background(), "42", "customer-1")
	require.NoError(t, err)
	assert.True(t, fetched)
	assert.True(t, liked)
	assert.False(t, published)
}

func TestPostService_Unlike_Succeeds(t *testing.T) {
	t.Parallel()

	unliked := false
	now := time.Now()

	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "author-1", CreatedAt: now}, nil
		},
		unlikeFn: func(ctx context.Context, postID string, customerID string) error {
			unliked = true
			return nil
		},
	}

	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	err := svc.Unlike(context.Background(), "42", "customer-1")
	require.NoError(t, err)
	assert.True(t, unliked)
}

func TestPostService_Like_PropagatesError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("like failed")
	now := time.Now()

	repo := &postRepositoryMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "author-1", CreatedAt: now}, nil
		},
		likeFn: func(ctx context.Context, postID string, customerID string) error {
			return repoErr
		},
	}

	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	err := svc.Like(context.Background(), "42", "customer-1")
	require.ErrorIs(t, err, repoErr)
}
