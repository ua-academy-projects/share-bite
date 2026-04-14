package collection

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

func TestCreateCollection(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()

		name        = gofakeit.ProductName()
		description = gofakeit.ProductDescription()

		errRepo = errors.New("unexpected database error")
	)

	tests := []struct {
		name string

		input entity.CreateCollectionInput

		mockFn func(repo *mockCollectionRepository)

		wantCollection entity.Collection
		wantErr        error
	}{
		{
			name: "success",
			input: entity.CreateCollectionInput{
				CustomerID:  customerID,
				Name:        name,
				Description: &description,
				IsPublic:    true,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("CreateCollection", mock.Anything, entity.CreateCollectionInput{
					CustomerID:  customerID,
					Name:        name,
					Description: &description,
					IsPublic:    true,
				}).
					Return(
						entity.Collection{
							ID:          collectionID,
							CustomerID:  customerID,
							Name:        name,
							Description: &description,
							IsPublic:    true,
						},
						nil,
					).Once()
			},
			wantCollection: entity.Collection{
				ID:          collectionID,
				CustomerID:  customerID,
				Name:        name,
				Description: &description,
				IsPublic:    true,
			},
			wantErr: nil,
		},
		{
			name: "error - repository fails",
			input: entity.CreateCollectionInput{
				CustomerID:  customerID,
				Name:        name,
				Description: &description,
				IsPublic:    true,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("CreateCollection", mock.Anything, entity.CreateCollectionInput{
					CustomerID:  customerID,
					Name:        name,
					Description: &description,
					IsPublic:    true,
				},
				).
					Return(entity.Collection{}, errRepo).Once()
			},
			wantCollection: entity.Collection{},
			wantErr:        errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := new(mockCollectionRepository)
			txManager := new(mockTxManager)
			businessClient := new(mockBusinessClient)
			svc := New(repo, txManager, businessClient)

			tt.mockFn(repo)

			collection, err := svc.CreateCollection(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantCollection, collection)
			repo.AssertExpectations(t)
		})
	}
}
