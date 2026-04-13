package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func strPtr(str string) *string { return &str }

func TestUpdate(t *testing.T) {
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

		input  entity.UpdateCustomer
		mockFn func(repo *mockCustomerRepository)

		wantCustomer entity.Customer
		wantErr      error
	}{
		{
			name: "success",
			input: entity.UpdateCustomer{
				UserID:          userID,
				UserName:        &userName,
				FirstName:       &firstName,
				LastName:        &lastName,
				AvatarObjectKey: &avatarObjectKey,
				Bio:             &bio,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{
						ID:              customerID,
						UserID:          userID,
						UserName:        gofakeit.Username(),
						FirstName:       gofakeit.FirstName(),
						LastName:        gofakeit.LastName(),
						AvatarObjectKey: strPtr(gofakeit.Word()),
						Bio:             strPtr(gofakeit.Hobby()),
					}, nil).
					Once()

				repo.On("Update", mock.Anything, entity.UpdateCustomer{
					UserID:          userID,
					UserName:        &userName,
					FirstName:       &firstName,
					LastName:        &lastName,
					AvatarObjectKey: &avatarObjectKey,
					Bio:             &bio,
				}).
					Return(entity.Customer{
						ID:              customerID,
						UserID:          userID,
						UserName:        userName,
						FirstName:       firstName,
						LastName:        lastName,
						AvatarObjectKey: &avatarObjectKey,
						Bio:             &bio,
					}, nil).
					Once()
			},
			wantCustomer: entity.Customer{
				ID:              customerID,
				UserID:          userID,
				UserName:        userName,
				FirstName:       firstName,
				LastName:        lastName,
				AvatarObjectKey: &avatarObjectKey,
				Bio:             &bio,
			},
			wantErr: nil,
		},
		{
			name: "error - get by user id repo fails",
			input: entity.UpdateCustomer{
				UserID:          userID,
				UserName:        &userName,
				FirstName:       &firstName,
				LastName:        &lastName,
				AvatarObjectKey: &avatarObjectKey,
				Bio:             &bio,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{}, errRepo).
					Once()
			},
			wantCustomer: entity.Customer{},
			wantErr:      errRepo,
		},
		{
			name: "error - update repo fails",
			input: entity.UpdateCustomer{
				UserID:          userID,
				UserName:        &userName,
				FirstName:       &firstName,
				LastName:        &lastName,
				AvatarObjectKey: &avatarObjectKey,
				Bio:             &bio,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{
						ID:              customerID,
						UserID:          userID,
						UserName:        gofakeit.Username(),
						FirstName:       gofakeit.FirstName(),
						LastName:        gofakeit.LastName(),
						AvatarObjectKey: strPtr(gofakeit.Word()),
						Bio:             strPtr(gofakeit.Hobby()),
					}, nil).
					Once()

				repo.On("Update", mock.Anything, entity.UpdateCustomer{
					UserID:          userID,
					UserName:        &userName,
					FirstName:       &firstName,
					LastName:        &lastName,
					AvatarObjectKey: &avatarObjectKey,
					Bio:             &bio,
				}).
					Return(entity.Customer{}, errRepo).
					Once()
			},
			wantCustomer: entity.Customer{},
			wantErr:      errRepo,
		},
		{
			name: "error - empty update",
			input: entity.UpdateCustomer{
				UserID: userID,
			},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{
						ID:              customerID,
						UserID:          userID,
						UserName:        gofakeit.Username(),
						FirstName:       gofakeit.FirstName(),
						LastName:        gofakeit.LastName(),
						AvatarObjectKey: strPtr(gofakeit.Word()),
						Bio:             strPtr(gofakeit.Hobby()),
					}, nil).
					Once()

				repo.On("Update", mock.Anything, entity.UpdateCustomer{
					UserID: userID,
				}).
					Return(entity.Customer{}, apperror.ErrEmptyUpdate).
					Once()
			},
			wantCustomer: entity.Customer{},
			wantErr:      apperror.ErrEmptyUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := new(mockCustomerRepository)
			svc := New(repo)
			tt.mockFn(repo)

			updatedCustomer, err := svc.Update(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantCustomer, updatedCustomer)
			repo.AssertExpectations(t)
		})
	}
}
