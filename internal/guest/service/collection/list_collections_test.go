package collection

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func TestListCollections(t *testing.T) {
	t.Parallel()

	var (
		customerID = gofakeit.UUID()

		errRepo  = errors.New("unexpected database error")
		baseTime = time.Now().UTC().Truncate(time.Millisecond)
	)

	generateToken := func(createdAt time.Time, id string) string {
		timeStr := createdAt.Format(time.RFC3339Nano)
		raw := fmt.Sprintf("%s|%s", timeStr, id)
		return base64.URLEncoding.EncodeToString([]byte(raw))
	}

	col1 := entity.Collection{ID: gofakeit.UUID(), CustomerID: customerID, CreatedAt: baseTime}
	col2 := entity.Collection{ID: gofakeit.UUID(), CustomerID: customerID, CreatedAt: baseTime.Add(-time.Hour)}
	col3 := entity.Collection{ID: gofakeit.UUID(), CustomerID: customerID, CreatedAt: baseTime.Add(-2 * time.Hour)}

	tests := []struct {
		name string

		input entity.ListCustomerCollectionsInput

		mockFn func(repo *mockCollectionRepository)

		wantOutput entity.ListCustomerCollectionsOutput
		wantErr    error
	}{
		{
			name: "success - default limit, no next page",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				PageSize:   0,
				PageToken:  "",
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListCustomerCollections", mock.Anything, customerID, time.Time{}, "", defaultLimit+1).
					Return([]entity.Collection{col1, col2}, nil).Once()
			},
			wantOutput: entity.ListCustomerCollectionsOutput{
				Collections:   []entity.Collection{col1, col2},
				NextPageToken: "",
			},
			wantErr: nil,
		},
		{
			name: "success - exact limit, has next page",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				PageSize:   2,
				PageToken:  "",
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListCustomerCollections", mock.Anything, customerID, time.Time{}, "", 3).
					Return([]entity.Collection{col1, col2, col3}, nil).Once()
			},
			wantOutput: entity.ListCustomerCollectionsOutput{
				Collections:   []entity.Collection{col1, col2},
				NextPageToken: generateToken(col2.CreatedAt, col2.ID),
			},
			wantErr: nil,
		},
		{
			name: "success - with token, max limit cap",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				PageSize:   150,
				PageToken:  generateToken(col1.CreatedAt, col1.ID),
			},
			mockFn: func(repo *mockCollectionRepository) {
				repo.On("ListCustomerCollections", mock.Anything, customerID, col1.CreatedAt, col1.ID, 101).
					Return([]entity.Collection{col2}, nil).Once()
			},
			wantOutput: entity.ListCustomerCollectionsOutput{
				Collections:   []entity.Collection{col2},
				NextPageToken: "",
			},
			wantErr: nil,
		},
		{
			name: "error - invalid token",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				PageSize:   10,
				PageToken:  "invalid-token",
			},
			mockFn:     func(repo *mockCollectionRepository) {},
			wantOutput: entity.ListCustomerCollectionsOutput{},
			wantErr:    apperror.ErrInvalidPageToken,
		},
		{
			name: "error - repository fails",
			input: entity.ListCustomerCollectionsInput{
				CustomerID: customerID,
				PageSize:   10,
				PageToken:  "",
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
