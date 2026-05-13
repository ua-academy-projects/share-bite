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
	"github.com/ua-academy-projects/share-bite/pkg/outbox"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		customerID = gofakeit.UUID()
		userID     = gofakeit.UUID()

		email     = gofakeit.Email()
		userName  = gofakeit.Username()
		firstName = gofakeit.Person().FirstName
		lastName  = gofakeit.Person().LastName

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
			svc := New(repo, outboxWriter, txManager)

			txManager.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil)

			if tt.wantErr == nil {
				outboxWriter.On("Enqueue", mock.Anything, mock.MatchedBy(func(event outbox.Event) bool {
					message, ok := event.Payload.(outbox.Message)
					return ok &&
						event.EventType == outbox.EventTypeRegistrationConfirmed &&
						message.EventType == outbox.EventTypeRegistrationConfirmed &&
						message.RecipientID == tt.input.UserID &&
						message.Metadata["email"] == tt.input.Email &&
						message.Metadata["username"] == tt.input.UserName
				})).Return(nil).Once()
			}
			tt.mockFn(repo)

			createdID, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantID, createdID)
			repo.AssertExpectations(t)
			outboxWriter.AssertExpectations(t)
			txManager.AssertExpectations(t)
		})
	}
}
