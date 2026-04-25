package collection

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func TestListCollections(t *testing.T) {
	t.Parallel()

	var (
		customerID = gofakeit.UUID()

		errRepo  = errors.New("unexpected database error")
		baseTime = time.Now().UTC().Truncate(time.Millisecond)
	)

	col1 := entity.Collection{ID: gofakeit.UUID(), CustomerID: customerID, CreatedAt: baseTime}
	col2 := entity.Collection{ID: gofakeit.UUID(), CustomerID: customerID, CreatedAt: baseTime.Add(-time.Hour)}
	col3 := entity.Collection{ID: gofakeit.UUID(), CustomerID: customerID, CreatedAt: baseTime.Add(-2 * time.Hour)}

	col2Time := col2.CreatedAt
	col2ID := col2.ID

	tests := []struct {
		name string

		input entity.ListCustomerCollectionsInput

		mockFn func(repo *mockCollectionRepository)

		wantOutput entity.ListCustomerCollectionsOutput
		wantErr    error
	}{
		{
			name: "success - no next page (mapper passed default limit + 1)",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				CursorTime: time.Time{},
				CursorID:   "",
				Limit:      21,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListCustomerCollections", mock.Anything, customerID, time.Time{}, "", 21).
					Return([]entity.Collection{col1, col2}, nil).Once()
			},
			wantOutput: entity.ListCustomerCollectionsOutput{
				Collections:    []entity.Collection{col1, col2},
				NextCursorTime: nil,
				NextCursorID:   nil,
			},
			wantErr: nil,
		},
		{
			name: "success - has next page (mapper passed limit + 1)",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				CursorTime: time.Time{},
				CursorID:   "",
				Limit:      3,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListCustomerCollections", mock.Anything, customerID, time.Time{}, "", 3).
					Return([]entity.Collection{col1, col2, col3}, nil).Once()
			},
			wantOutput: entity.ListCustomerCollectionsOutput{
				Collections:    []entity.Collection{col1, col2},
				NextCursorTime: &col2Time,
				NextCursorID:   &col2ID,
			},
			wantErr: nil,
		},
		{
			name: "success - with cursors, max limit cap",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				CursorTime: col1.CreatedAt,
				CursorID:   col1.ID,
				Limit:      101,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListCustomerCollections", mock.Anything, customerID, col1.CreatedAt, col1.ID, 101).
					Return([]entity.Collection{col2}, nil).Once()
			},
			wantOutput: entity.ListCustomerCollectionsOutput{
				Collections:    []entity.Collection{col2},
				NextCursorTime: nil,
				NextCursorID:   nil,
			},
			wantErr: nil,
		},
		{
			name: "error - repository fails",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				CursorTime: time.Time{},
				CursorID:   "",
				Limit:      11,
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListCustomerCollections", mock.Anything, customerID, time.Time{}, "", 11).
					Return([]entity.Collection(nil), errRepo).Once()
			},
			wantOutput: entity.ListCustomerCollectionsOutput{},
			wantErr:    errRepo,
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

			output, err := svc.ListCustomerCollections(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantOutput, output)
			repo.AssertExpectations(t)
		})
	}
}
