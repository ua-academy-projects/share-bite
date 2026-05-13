package customer

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"testing"
)

func TestGetByIDs(t *testing.T) {
	t.Parallel()

	var (
		customerID1 = gofakeit.UUID()
		customerID2 = gofakeit.UUID()

		userID1 = gofakeit.UUID()
		userID2 = gofakeit.UUID()

		userName1 = gofakeit.Username()
		userName2 = gofakeit.Username()

		errRepo = errors.New("unexpected repository error")
	)

	tests := []struct {
		name string

		ids    []string
		mockFn func(repo *mockCustomerRepository)

		wantCustomers []entity.Customer
		wantErr       error
	}{
		{
			name: "success",
			ids:  []string{customerID1, customerID2},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByIDs", mock.Anything, []string{customerID1, customerID2}).
					Return([]entity.Customer{
						{
							ID:       customerID1,
							UserID:   userID1,
							UserName: userName1,
						},
						{
							ID:       customerID2,
							UserID:   userID2,
							UserName: userName2,
						},
					}, nil).
					Once()
			},
			wantCustomers: []entity.Customer{
				{
					ID:       customerID1,
					UserID:   userID1,
					UserName: userName1,
				},
				{
					ID:       customerID2,
					UserID:   userID2,
					UserName: userName2,
				},
			},
			wantErr: nil,
		},
		{
			name: "error - repo fails",
			ids:  []string{customerID1, customerID2},
			mockFn: func(repo *mockCustomerRepository) {
				repo.On("GetByIDs", mock.Anything, []string{customerID1, customerID2}).
					Return(nil, errRepo).
					Once()
			},
			wantCustomers: nil,
			wantErr:       errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockCustomerRepository)
			svc := New(repo, nil, nil)

			tt.mockFn(repo)

			customers, err := svc.GetByIDs(context.Background(), tt.ids)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantCustomers, customers)
			repo.AssertExpectations(t)
		})
	}
}
