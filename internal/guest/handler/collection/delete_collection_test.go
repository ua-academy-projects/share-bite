package collection

import (
	"errors"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestDeleteCollection(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()

		invalidCollectionID = "bad-uuid"
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
			name:         "success",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("DeleteCollection", mock.Anything, collectionID, customerID).
					Return(nil).
					Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
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
			name:         "customer id missing in context",
			collectionID: collectionID,
			customerID:   nil,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "customer id invalid type in context",
			collectionID: collectionID,
			customerID:   123,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service returns forbidden",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("DeleteCollection", mock.Anything, collectionID, customerID).
					Return(apperror.ErrCollectionAccessDenied).
					Once()
			},
			wantCode: http.StatusForbidden,
			wantBody: response.ErrorResponse{Message: apperror.ErrCollectionAccessDenied.Error()},
		},
		{
			name:         "service returns not found",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				err := apperror.CollectionNotFoundID(collectionID)
				s.On("DeleteCollection", mock.Anything, collectionID, customerID).
					Return(err).
					Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: response.ErrorResponse{Message: apperror.CollectionNotFoundID(collectionID).Error()},
		},
		{
			name:         "service returns unknown error",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("DeleteCollection", mock.Anything, collectionID, customerID).
					Return(errors.New("db down")).
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

			r.DELETE("/collections/:collectionId", withCustomerID(tt.customerID, h.deleteCollection))

			w := performRequest(r, http.MethodDelete, "/collections/"+tt.collectionID)

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
