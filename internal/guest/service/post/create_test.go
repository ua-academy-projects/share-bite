package post

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

func TestPostService_Create_SucceedsWithoutImages(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		createFn: func(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
			return entity.Post{
				ID:         "42",
				CustomerID: in.CustomerID,
				VenueID:    in.VenueID,
				Text:       in.Text,
				Rating:     in.Rating,
				Status:     entity.PostStatusDraft,
				CreatedAt:  time.Now(),
			}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	out, err := svc.Create(context.Background(), dto.CreatePostInput{
		CustomerID: "customer-1",
		VenueID:    123,
		Text:       "hello",
		Rating:     5,
	})
	require.NoError(t, err)
	assert.Equal(t, "42", out.ID)
	assert.Equal(t, "customer-1", repo.lastCreateInput.CustomerID)
	assert.Equal(t, int64(123), repo.lastCreateInput.VenueID)
}

func TestPostService_Create_RejectsWhenVenueDoesNotExist(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{}
	svc := New(repo, &venueProviderMock{
		checkExistsFn: func(ctx context.Context, venueID int64) (bool, error) {
			return false, nil
		},
	}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	_, err := svc.Create(context.Background(), dto.CreatePostInput{VenueID: 777})
	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code.NotFound, appErr.Code)
}

func TestPostService_Create_WithMentions_Success(t *testing.T) {
	t.Parallel()

	called := false

	repo := &postRepositoryMock{
		createFn: func(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
			return entity.Post{ID: "post-1"}, nil
		},
		createMentionsFn: func(ctx context.Context, mentions []entity.PostMention) error {
			called = true
			require.Len(t, mentions, 2)
			return nil
		},
	}

	svc := New(
		repo,
		&venueProviderMock{},
		&followRepoMock{
			getAllowedMentionsFn: func(ctx context.Context, customerID string, ids []string) ([]string, error) {
				return ids, nil
			},
		},
		&customerRepoMock{},
		&storageMock{},
		&txManagerMock{},
	)

	_, err := svc.Create(context.Background(), dto.CreatePostInput{
		CustomerID: "user-1",
		VenueID:    1,
		Mentions:   []string{"user-2", "user-3"},
	})

	require.NoError(t, err)
	assert.True(t, called)
}

func TestPostService_Create_MentionYourself_Error(t *testing.T) {
	t.Parallel()

	svc := New(
		&postRepositoryMock{},
		&venueProviderMock{},
		&followRepoMock{},
		&customerRepoMock{},
		&storageMock{},
		&txManagerMock{},
	)

	_, err := svc.Create(context.Background(), dto.CreatePostInput{
		CustomerID: "user-1",
		VenueID:    1,
		Mentions:   []string{"user-1"},
	})

	require.Error(t, err)
}

func TestPostService_Create_TooManyMentions_Error(t *testing.T) {
	t.Parallel()

	mentions := make([]string, 11)
	for i := range mentions {
		mentions[i] = fmt.Sprintf("user-%d", i)
	}

	svc := New(
		&postRepositoryMock{},
		&venueProviderMock{},
		&followRepoMock{},
		&customerRepoMock{},
		&storageMock{},
		&txManagerMock{},
	)

	_, err := svc.Create(context.Background(), dto.CreatePostInput{
		CustomerID: "user-1",
		VenueID:    1,
		Mentions:   mentions,
	})

	require.Error(t, err)
}

func TestPostService_Create_ForbiddenMention_Error(t *testing.T) {
	t.Parallel()

	svc := New(
		&postRepositoryMock{},
		&venueProviderMock{},
		&followRepoMock{
			getAllowedMentionsFn: func(ctx context.Context, customerID string, ids []string) ([]string, error) {
				return []string{}, nil
			},
		},
		nil,
		&storageMock{},
		&txManagerMock{},
	)

	_, err := svc.Create(context.Background(), dto.CreatePostInput{
		CustomerID: "user-1",
		VenueID:    1,
		Mentions:   []string{"user-2"},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, apperror.ErrForbiddenMention)
}
