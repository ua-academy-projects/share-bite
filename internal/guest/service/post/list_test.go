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
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

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
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	_, err := svc.List(context.Background(), dto.ListPostsInput{})
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}

func TestPostService_List_WithMentions(t *testing.T) {
	t.Parallel()

	now := time.Now()

	repo := &postRepositoryMock{
		listFn: func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
			return dto.ListPostsOutput{
				Posts: []entity.Post{
					{ID: "1", CreatedAt: now},
				},
				Total: 1,
			}, nil
		},
		listMentionsFn: func(ctx context.Context, postIDs []string) (map[string][]entity.PostMention, error) {
			return map[string][]entity.PostMention{
				"1": {
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

	out, err := svc.List(context.Background(), dto.ListPostsInput{})
	require.NoError(t, err)

	require.Len(t, out.Posts, 1)
	require.Len(t, out.Posts[0].Mentions, 2)

	assert.Equal(t, "user-1", out.Posts[0].Mentions[0].ID)
	assert.Equal(t, "user-2", out.Posts[0].Mentions[1].ID)
}

func TestPostService_List_NoMentions(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		listFn: func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
			return dto.ListPostsOutput{
				Posts: []entity.Post{
					{ID: "1"},
				},
				Total: 1,
			}, nil
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

	out, err := svc.List(context.Background(), dto.ListPostsInput{})
	require.NoError(t, err)

	assert.Empty(t, out.Posts[0].Mentions)
}

func TestPostService_List_GetByIDs_Error(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		listFn: func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
			return dto.ListPostsOutput{
				Posts: []entity.Post{{ID: "1"}},
			}, nil
		},
		listMentionsFn: func(ctx context.Context, postIDs []string) (map[string][]entity.PostMention, error) {
			return map[string][]entity.PostMention{
				"1": {
					{CustomerID: "user-1"},
				},
			}, nil
		},
	}

	customerRepo := &customerRepoMock{
		getByIDsFn: func(ctx context.Context, ids []string) ([]entity.Customer, error) {
			return nil, errors.New("boom")
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

	_, err := svc.List(context.Background(), dto.ListPostsInput{})
	require.Error(t, err)
}
