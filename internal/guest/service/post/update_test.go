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

func TestPostService_Update_AllowsDraftToPublished(t *testing.T) {
	t.Parallel()

	now := time.Now()
	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusDraft, CreatedAt: now}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	status := entity.PostStatusPublished
	_, err := svc.Update(context.Background(), entity.UpdatePostInput{
		ID:         "42",
		CustomerID: "customer-1",
		Status:     &status,
	})
	require.NoError(t, err)
	require.NotNil(t, repo.lastUpdateInput.Status)
	assert.Equal(t, entity.PostStatusPublished, *repo.lastUpdateInput.Status)
}

func TestPostService_Update_RejectsPublishedToDraft(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusPublished}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	status := entity.PostStatusDraft
	_, err := svc.Update(context.Background(), entity.UpdatePostInput{
		ID:         "42",
		CustomerID: "customer-1",
		Status:     &status,
	})
	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code.InvalidRequest, appErr.Code)
}

func TestPostService_Update_AllowsEditingArchivedPost(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusArchived}, nil
		},
		updateFn: func(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
			return entity.Post{ID: in.ID, CustomerID: in.CustomerID, Text: *in.Text}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	text := "updated text"
	_, err := svc.Update(context.Background(), entity.UpdatePostInput{
		ID:         "42",
		CustomerID: "customer-1",
		Text:       &text,
	})
	require.NoError(t, err)
	require.NotNil(t, repo.lastUpdateInput.Text)
	assert.Equal(t, "updated text", *repo.lastUpdateInput.Text)
}

func TestPostService_Update_AllowsUnarchiveArchivedToPublished(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusArchived}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	status := entity.PostStatusPublished
	_, err := svc.Update(context.Background(), entity.UpdatePostInput{
		ID:         "42",
		CustomerID: "customer-1",
		Status:     &status,
	})
	require.NoError(t, err)
	require.NotNil(t, repo.lastUpdateInput.Status)
	assert.Equal(t, entity.PostStatusPublished, *repo.lastUpdateInput.Status)
}

func TestPostService_Update_RejectsForeignOwner(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusDraft}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	status := entity.PostStatusPublished
	_, err := svc.Update(context.Background(), entity.UpdatePostInput{
		ID:         "42",
		CustomerID: "customer-2",
		Status:     &status,
	})
	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code.NotFound, appErr.Code)
}

func TestPostService_Update_RejectsDraftToDeleted(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(ctx context.Context, postID string) (entity.Post, error) {
			return entity.Post{ID: postID, CustomerID: "customer-1", Status: entity.PostStatusDraft}, nil
		},
	}
	svc := New(repo, &venueProviderMock{}, &followRepoMock{}, &customerRepoMock{}, &storageMock{}, &txManagerMock{})

	status := entity.PostStatusDeleted
	_, err := svc.Update(context.Background(), entity.UpdatePostInput{
		ID:         "42",
		CustomerID: "customer-1",
		Status:     &status,
	})
	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, code.InvalidRequest, appErr.Code)
}

func TestPostService_Update_AllowsAcceptedCollaboratorToEditText(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(
			ctx context.Context,
			postID string,
		) (entity.Post, error) {
			return entity.Post{
				ID:         postID,
				CustomerID: "owner-1",
				Status:     entity.PostStatusPublished,
			}, nil
		},
		isAcceptedCollaboratorFn: func(
			ctx context.Context,
			postID string,
			customerID string,
		) (bool, error) {
			return true, nil
		},
		updateFn: func(
			ctx context.Context,
			in entity.UpdatePostInput,
		) (entity.Post, error) {
			return entity.Post{
				ID:         in.ID,
				CustomerID: "owner-1",
				Text:       *in.Text,
			}, nil
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

	text := "updated by collaborator"

	_, err := svc.Update(
		context.Background(),
		entity.UpdatePostInput{
			ID:         "42",
			CustomerID: "collab-1",
			Text:       &text,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, repo.lastUpdateInput.Text)

	assert.Equal(
		t,
		"updated by collaborator",
		*repo.lastUpdateInput.Text,
	)
}

func TestPostService_Update_RejectsCollaboratorStatusChange(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(
			ctx context.Context,
			postID string,
		) (entity.Post, error) {
			return entity.Post{
				ID:         postID,
				CustomerID: "owner-1",
				Status:     entity.PostStatusDraft,
			}, nil
		},
		isAcceptedCollaboratorFn: func(
			ctx context.Context,
			postID string,
			customerID string,
		) (bool, error) {
			return true, nil
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

	status := entity.PostStatusPublished

	_, err := svc.Update(
		context.Background(),
		entity.UpdatePostInput{
			ID:         "42",
			CustomerID: "collab-1",
			Status:     &status,
		},
	)

	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)

	assert.Equal(t, code.Forbidden, appErr.Code)
}

func TestPostService_Update_RejectsCollaboratorRatingChange(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(
			ctx context.Context,
			postID string,
		) (entity.Post, error) {
			return entity.Post{
				ID:         postID,
				CustomerID: "owner-1",
				Status:     entity.PostStatusPublished,
			}, nil
		},
		isAcceptedCollaboratorFn: func(
			ctx context.Context,
			postID string,
			customerID string,
		) (bool, error) {
			return true, nil
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

	rating := int16(5)

	_, err := svc.Update(
		context.Background(),
		entity.UpdatePostInput{
			ID:         "42",
			CustomerID: "collab-1",
			Rating:     &rating,
		},
	)

	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)

	assert.Equal(t, code.Forbidden, appErr.Code)
}

func TestPostService_Update_RejectsCollaboratorVenueChange(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		getByIDFn: func(
			ctx context.Context,
			postID string,
		) (entity.Post, error) {
			return entity.Post{
				ID:         postID,
				CustomerID: "owner-1",
				Status:     entity.PostStatusPublished,
			}, nil
		},
		isAcceptedCollaboratorFn: func(
			ctx context.Context,
			postID string,
			customerID string,
		) (bool, error) {
			return true, nil
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

	venueID := int64(123)

	_, err := svc.Update(
		context.Background(),
		entity.UpdatePostInput{
			ID:         "42",
			CustomerID: "collab-1",
			VenueID:    &venueID,
		},
	)

	require.Error(t, err)

	var appErr *apperror.Error
	require.ErrorAs(t, err, &appErr)

	assert.Equal(t, code.Forbidden, appErr.Code)
}
