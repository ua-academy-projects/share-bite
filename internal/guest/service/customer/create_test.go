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

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		customerID = gofakeit.UUID()
		userID     = gofakeit.UUID()

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
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
				Bio:       &bio,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("Create", mock.Anything, entity.CreateCustomer{
					UserID:    userID,
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
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
				Bio:       &bio,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("Create", mock.Anything, entity.CreateCustomer{
					UserID:    userID,
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
			svc := New(repo)
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
		})
	}
}
