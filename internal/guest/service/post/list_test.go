package post

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func TestPostService_List_Succeeds(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &postRepositoryMock{
		listFn: func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
			return dto.ListPostsOutput{
				Posts: []entity.Post{{ID: "1", CreatedAt: now}},
				Total: 1,
			}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	out, err := svc.List(context.Background(), dto.ListPostsInput{Limit: 10, Offset: 5, CustomerID: "c-1"})
	require.NoError(t, err)
	assert.Equal(t, 1, out.Total)
	require.Len(t, out.Posts, 1)
	assert.Equal(t, "1", out.Posts[0].ID)
	assert.Equal(t, now, out.Posts[0].CreatedAt)
	assert.Equal(t, 10, repo.lastListInput.Limit)
	assert.Equal(t, 5, repo.lastListInput.Offset)
	assert.Equal(t, "c-1", repo.lastListInput.CustomerID)
}

func TestPostService_List_PropagatesRepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo list failed")
	repo := &postRepositoryMock{
		listFn: func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
			return dto.ListPostsOutput{}, repoErr
		},
	}
	svc := New(repo, &venueProviderMock{}, &storageMock{}, &txManagerMock{})

	_, err := svc.List(context.Background(), dto.ListPostsInput{})
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
