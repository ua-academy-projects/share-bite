package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

func TestCreate_EmailResolutionFailures(t *testing.T) {
	t.Parallel()

	var (
		userID    = gofakeit.UUID()
		email     = gofakeit.Email()
		userName  = gofakeit.Username()
		firstName = gofakeit.Person().FirstName
		lastName  = gofakeit.Person().LastName
	)

	tests := []struct {
		name       string
		ctx        context.Context
		setupAdmin func(adminClient *mockEmailClient, input entity.CreateCustomer)
		wantErrMsg string
	}{
		{
			name:       "missing access token",
			ctx:        context.Background(),
			wantErrMsg: "missing access token to resolve customer email",
		},
		{
			name: "admin email lookup fails",
			ctx:  context.WithValue(context.Background(), middleware.CtxAccessToken, "access-token"),
			setupAdmin: func(adminClient *mockEmailClient, input entity.CreateCustomer) {
				emailErr := errors.New("admin email lookup failed")
				adminClient.On("GetUserEmail", mock.Anything, input.UserID, "access-token").Return("", emailErr).Once()
			},
			wantErrMsg: "admin email lookup failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockCustomerRepository)
			outboxWriter := new(mockOutboxWriter)
			txManager := new(mockTxManager)
			adminClient := new(mockEmailClient)
			svc := New(repo, outboxWriter, txManager, adminClient)

			input := entity.CreateCustomer{
				UserID:    userID,
				Email:     email,
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
			}

			if tt.setupAdmin != nil {
				tt.setupAdmin(adminClient, input)
			}

			_, err := svc.Create(tt.ctx, input)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrMsg)

			repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			outboxWriter.AssertNotCalled(t, "Enqueue", mock.Anything, mock.Anything)
			txManager.AssertNotCalled(t, "ReadCommitted", mock.Anything, mock.Anything)
			adminClient.AssertExpectations(t)
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		customerID = gofakeit.UUID()
		userID     = gofakeit.UUID()

		email      = gofakeit.Email()
		adminEmail = gofakeit.Email()
		userName   = gofakeit.Username()
		firstName  = gofakeit.Person().FirstName
		lastName   = gofakeit.Person().LastName

		bio = gofakeit.Person().Hobby

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name string

		input  entity.CreateCustomer
		mockFn func(repo *mockCustomerRepository)

		wantID  string
		wantErr error
	}{
		{
			name: "success",
			input: entity.CreateCustomer{
				UserID:    userID,
				Email:     email,
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
				Bio:       &bio,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("Create", mock.Anything, entity.CreateCustomer{
					UserID:    userID,
					Email:     email,
					UserName:  userName,
					FirstName: firstName,
					LastName:  lastName,
					Bio:       &bio,
				}).Return(customerID, nil).Once()
			},
			wantID:  customerID,
			wantErr: nil,
		},
		{
			name: "error - create repo fails",
			input: entity.CreateCustomer{
				UserID:    userID,
				Email:     email,
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
				Bio:       &bio,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("Create", mock.Anything, entity.CreateCustomer{
					UserID:    userID,
					Email:     email,
					UserName:  userName,
					FirstName: firstName,
					LastName:  lastName,
					Bio:       &bio,
				}).Return("", errRepo).Once()
			},
			wantID:  "",
			wantErr: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := new(mockCustomerRepository)
			outboxWriter := new(mockOutboxWriter)
			txManager := new(mockTxManager)
			adminClient := new(mockEmailClient)
			svc := New(repo, outboxWriter, txManager, adminClient)

			ctx := context.WithValue(context.Background(), middleware.CtxAccessToken, "access-token")

			txManager.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil)

			adminClient.On("GetUserEmail", mock.Anything, tt.input.UserID, "access-token").Return(adminEmail, nil).Once()

			if tt.wantErr == nil {
				outboxWriter.On("Enqueue", mock.Anything, mock.MatchedBy(func(event outbox.Event) bool {
					message, ok := event.Payload.(outbox.Message)
					return ok &&
						event.EventType == outbox.EventTypeRegistrationConfirmed &&
						message.EventType == outbox.EventTypeRegistrationConfirmed &&
						message.RecipientID == tt.input.UserID &&
						message.Metadata["email"] == adminEmail &&
						message.Metadata["username"] == tt.input.UserName
				})).Return(nil).Once()
			}
			tt.mockFn(repo)

			createdID, err := svc.Create(ctx, tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantID, createdID)
			repo.AssertExpectations(t)
			adminClient.AssertExpectations(t)
			outboxWriter.AssertExpectations(t)
			txManager.AssertExpectations(t)
		})
	}
}
