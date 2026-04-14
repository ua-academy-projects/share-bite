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
)

func TestGetByUserID(t *testing.T) {
	t.Parallel()

	var (
		userID = gofakeit.UUID()

		customerID = gofakeit.UUID()
		userName   = gofakeit.Username()
		firstName  = gofakeit.Person().FirstName
		lastName   = gofakeit.Person().LastName

		bio             = gofakeit.Person().Hobby
		avatarObjectKey = gofakeit.Word()

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name string

		userID string
		mockFn func(repo *mockCustomerRepository)

		wantCustomer entity.Customer
		wantErr      error
	}{
		{
			name:   "success",
			userID: userID,
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{
						ID:              customerID,
						UserID:          userID,
						UserName:        userName,
						FirstName:       firstName,
						LastName:        lastName,
						Bio:             &bio,
						AvatarObjectKey: &avatarObjectKey,
					}, nil).
					Once()
			},
			wantCustomer: entity.Customer{
				ID:              customerID,
				UserID:          userID,
				UserName:        userName,
				FirstName:       firstName,
				LastName:        lastName,
				Bio:             &bio,
				AvatarObjectKey: &avatarObjectKey,
			},
			wantErr: nil,
		},
		{
			name:   "error - get by user id repo fails",
			userID: userID,
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{}, errRepo).
					Once()
			},
			wantCustomer: entity.Customer{},
			wantErr:      errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := new(mockCustomerRepository)
			svc := New(repo)
			tt.mockFn(repo)

			customer, err := svc.GetByUserID(context.Background(), tt.userID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantCustomer, customer)
			repo.AssertExpectations(t)
		})
	}
}
