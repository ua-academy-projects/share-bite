package collection

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestGetCollection(t *testing.T) {
	t.Parallel()

	var (
		collectionID        = gofakeit.UUID()
		customerID          = gofakeit.UUID()
		invalidCollectionID = "bad-uuid"

		now = time.Now().UTC()
	)

	tests := []struct {
		name string

		collectionID string
		customerID   any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success with customer id",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				cid := customerID
				s.On("GetCollection", mock.Anything, collectionID, &cid).
					Return(entity.Collection{
						ID:          collectionID,
						CustomerID:  customerID,
						Name:        "My collection",
						Description: strPtr("desc"),
						IsPublic:    false,
						CreatedAt:   now,
						UpdatedAt:   now,
					}, nil).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: getCollectionResponse{
				Collection: collectionResponse{
					ID:          collectionID,
					Name:        "My collection",
					Description: strPtr("desc"),
					IsPublic:    false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:         "success without customer id (public access path)",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				s.On("GetCollection", mock.Anything, collectionID, (*string)(nil)).
					Return(entity.Collection{
						ID:         collectionID,
						CustomerID: gofakeit.UUID(),
						Name:       "Public",
						IsPublic:   true,
						CreatedAt:  now,
						UpdatedAt:  now,
					}, nil).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: getCollectionResponse{
				Collection: collectionResponse{
					ID:        collectionID,
					Name:      "Public",
					IsPublic:  true,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
		{
			name:         "validation error - invalid collectionId uuid",
			collectionID: invalidCollectionID,
			customerID:   customerID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "collectionId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:         "optional customer id invalid type in context",
			collectionID: collectionID,
			customerID:   123,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service returns not found",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				err := apperror.CollectionNotFoundID(collectionID)
				s.On("GetCollection", mock.Anything, collectionID, (*string)(nil)).
					Return(entity.Collection{}, err).
					Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: response.ErrorResponse{Message: apperror.CollectionNotFoundID(collectionID).Error()},
		},
		{
			name:         "service returns unknown error",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				s.On("GetCollection", mock.Anything, collectionID, (*string)(nil)).
					Return(entity.Collection{}, errors.New("boom")).
					Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: response.ErrorResponse{Message: internalErrMsg},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := new(mockCollectionService)
			tt.mockFn(svc)

			h := &handler{service: svc}
			r := newTestRouter()

			r.GET("/collections/:collectionId", withCustomerID(tt.customerID, h.getCollection))

			w := performRequest(r, http.MethodGet, "/collections/"+tt.collectionID)

			require.Equal(t, tt.wantCode, w.Code)
			if tt.wantBody != nil {
				assertJSONBody(t, tt.wantBody, w.Body.String())
			} else {
				require.Empty(t, w.Body.String())
			}

			svc.AssertExpectations(t)
		})
	}
}
