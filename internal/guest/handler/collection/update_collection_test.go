package collection

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestUpdateCollection(t *testing.T) {
	t.Parallel()

	var (
		collectionID        = gofakeit.UUID()
		invalidCollectionID = "bad-uuid"
		customerID          = gofakeit.UUID()

		now = time.Now().UTC()
	)

	tests := []struct {
		name string

		collectionID string
		customerID   any
		body         any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success",
			collectionID: collectionID,
			customerID:   customerID,
			body: updateCollectionBody{
				Name:        strPtr("  The best collection ever  "),
				Description: strPtr("  I promise  "),
				IsPublic:    boolPtr(true),
			},
			mockFn: func(s *mockCollectionService) {
				name := "The best collection ever"
				desc := "I promise"
				isPublic := true

				s.On("UpdateCollection", mock.Anything, entity.UpdateCollectionInput{
					CollectionID: collectionID,
					CustomerID:   customerID,
					Name:         &name,
					Description:  &desc,
					IsPublic:     &isPublic,
				}).Return(entity.Collection{
					ID:          collectionID,
					CustomerID:  customerID,
					Name:        name,
					Description: &desc,
					IsPublic:    isPublic,
					CreatedAt:   now,
					UpdatedAt:   now,
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: updateCollectionResponse{
				Collection: collectionResponse{
					ID:          collectionID,
					Name:        "The best collection ever",
					Description: strPtr("I promise"),
					IsPublic:    true,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:         "validation error uri",
			collectionID: invalidCollectionID,
			customerID:   customerID,
			body:         updateCollectionBody{Name: strPtr("Not the best but it's cool")},
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
			name:         "validation error body",
			collectionID: collectionID,
			customerID:   customerID,
			body:         updateCollectionBody{Name: strPtr("")},
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "name", Message: "This field must be at least 1 characters long"},
				},
			},
		},
		{
			name:         "empty update",
			collectionID: collectionID,
			customerID:   customerID,
			body:         updateCollectionBody{},
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody:     response.ErrorResponse{Message: apperror.ErrEmptyUpdate.Error()},
		},
		{
			name:         "missing customer id in ctx",
			collectionID: collectionID,
			customerID:   nil,
			body:         updateCollectionBody{Name: strPtr("Good")},
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service forbidden",
			collectionID: collectionID,
			customerID:   customerID,
			body:         updateCollectionBody{Name: strPtr("Nice")},
			mockFn: func(s *mockCollectionService) {
				s.On("UpdateCollection", mock.Anything, mock.Anything).
					Return(entity.Collection{}, apperror.ErrCollectionAccessDenied).
					Once()
			},
			wantCode: http.StatusForbidden,
			wantBody: response.ErrorResponse{Message: apperror.ErrCollectionAccessDenied.Error()},
		},
		{
			name:         "service unknown error",
			collectionID: collectionID,
			customerID:   customerID,
			body:         updateCollectionBody{Name: strPtr("Not bad collection")},
			mockFn: func(s *mockCollectionService) {
				s.On("UpdateCollection", mock.Anything, mock.Anything).
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

			r.PATCH("/collections/:collectionId", withCustomerID(tt.customerID, h.updateCollection))

			var w *httptest.ResponseRecorder
			if s, ok := tt.body.(string); ok {
				w = performRawJSONRequest(t, r, http.MethodPatch, "/collections/"+tt.collectionID, s)
			} else {
				w = performJSONRequest(t, r, http.MethodPatch, "/collections/"+tt.collectionID, tt.body)
			}

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
