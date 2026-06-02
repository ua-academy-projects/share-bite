package post

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
	"testing"
)

func TestPostService_AcceptInvitation_Success(t *testing.T) {
	t.Parallel()

	accepted := false
	published := false

	repo := &postRepositoryMock{
		acceptPostInvitationFn: func(
			ctx context.Context,
			collaboratorID string,
			customerID string,
		) (string, error) {
			accepted = true
			return "post-1", nil
		},
		tryPublishPostIfAllAcceptedFn: func(
			ctx context.Context,
			postID string,
		) (bool, error) {
			published = true
			return false, nil
		},
	}

	svc := New(
		repo,
		&venueProviderMock{},
		&followRepoMock{},
		&customerRepoMock{},
		&storageMock{},
		&txManagerMock{},
		WithOutboxWriter(&outboxWriterMock{}),
	)

	err := svc.AcceptInvitation(
		context.Background(),
		"collab-1",
		"user-1",
	)

	require.NoError(t, err)
	assert.True(t, accepted)
	assert.True(t, published)
}

func TestPostService_AcceptInvitation_ReturnsErrorWhenAcceptFails(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		acceptPostInvitationFn: func(
			ctx context.Context,
			collaboratorID string,
			customerID string,
		) (string, error) {
			return "", assert.AnError
		},
	}

	svc := New(
		repo,
		&venueProviderMock{},
		&followRepoMock{},
		&customerRepoMock{},
		&storageMock{},
		&txManagerMock{},
		WithOutboxWriter(&outboxWriterMock{}),
	)

	err := svc.AcceptInvitation(
		context.Background(),
		"collab-1",
		"user-1",
	)

	require.Error(t, err)
	assert.ErrorContains(t, err, "accept invitation")
}

func TestPostService_AcceptInvitation_ReturnsErrorWhenPublishCheckFails(t *testing.T) {
	t.Parallel()

	repo := &postRepositoryMock{
		acceptPostInvitationFn: func(
			ctx context.Context,
			collaboratorID string,
			customerID string,
		) (string, error) {
			return "post-1", nil
		},
		tryPublishPostIfAllAcceptedFn: func(
			ctx context.Context,
			postID string,
		) (bool, error) {
			return false, assert.AnError
		},
	}

	svc := New(
		repo,
		&venueProviderMock{},
		&followRepoMock{},
		&customerRepoMock{},
		&storageMock{},
		&txManagerMock{},
		WithOutboxWriter(&outboxWriterMock{}),
	)

	err := svc.AcceptInvitation(
		context.Background(),
		"collab-1",
		"user-1",
	)

	require.Error(t, err)
	assert.ErrorContains(t, err, "try publish post")
}

func TestPostService_AcceptInvitation_PublishesNotificationsWhenAllAccepted(t *testing.T) {
	t.Parallel()

	enqueued := 0

	repo := &postRepositoryMock{
		acceptPostInvitationFn: func(
			ctx context.Context,
			collaboratorID string,
			customerID string,
		) (string, error) {
			return "post-1", nil
		},

		tryPublishPostIfAllAcceptedFn: func(
			ctx context.Context,
			postID string,
		) (bool, error) {
			return true, nil
		},

		getAcceptedPostCollaboratorsFn: func(
			ctx context.Context,
			postID string,
		) ([]string, error) {
			return []string{
				"user-2",
				"user-3",
			}, nil
		},

		getAuthorCustomerIDFn: func(
			ctx context.Context,
			postID string,
		) (string, error) {
			return "author-1", nil
		},
	}

	customerRepo := &customerRepoMock{
		getByIDFn: func(
			ctx context.Context,
			customerID string,
		) (entity.Customer, error) {
			return entity.Customer{
				ID:     customerID,
				UserID: customerID + "-user",
			}, nil
		},
	}

	outboxWriter := &outboxWriterMock{
		enqueueFn: func(
			ctx context.Context,
			event outbox.Event,
		) error {

			enqueued++

			assert.Equal(
				t,
				outbox.EventTypePostPublished,
				event.EventType,
			)

			return nil
		},
	}

	svc := New(
		repo,
		&venueProviderMock{},
		&followRepoMock{},
		customerRepo,
		&storageMock{},
		&txManagerMock{},
		WithOutboxWriter(outboxWriter),
	)

	err := svc.AcceptInvitation(
		context.Background(),
		"collab-1",
		"user-1",
	)

	require.NoError(t, err)

	assert.Equal(t, 3, enqueued)
}
