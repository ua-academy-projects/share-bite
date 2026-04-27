package collection

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestAddVenue(t *testing.T) {
	t.Parallel()

	var (
		collectionID = gofakeit.UUID()
		customerID   = gofakeit.UUID()
		venueID      = int64(gofakeit.Number(1, 1_000_000))
		venueIDStr   = strconv.FormatInt(venueID, 10)

		invalidCollectionID = "bad-uuid"
		invalidVenueIDStr   = "0"
	)

	tests := []struct {
		name string

		collectionID string
		venueID      string
		customerID   any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success",
			collectionID: collectionID,
			venueID:      venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("AddVenue", mock.Anything, collectionID, customerID, venueID).
					Return(nil).
					Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:         "validation error - invalid collectionId uuid",
			collectionID: invalidCollectionID,
			venueID:      venueIDStr,
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
			name:         "validation error - venueId must be gte 1",
			collectionID: collectionID,
			venueID:      invalidVenueIDStr,
			customerID:   customerID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "venueId", Message: "This field is required"},
				},
			},
		},
		{
			name:         "customer id missing in context",
			collectionID: collectionID,
			venueID:      venueIDStr,
			customerID:   nil,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "customer id invalid type in context",
			collectionID: collectionID,
			venueID:      venueIDStr,
			customerID:   123,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service returns forbidden",
			collectionID: collectionID,
			venueID:      venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("AddVenue", mock.Anything, collectionID, customerID, venueID).
					Return(apperror.ErrCollectionAccessDenied).
					Once()
			},
			wantCode: http.StatusForbidden,
			wantBody: response.ErrorResponse{Message: apperror.ErrCollectionAccessDenied.Error()},
		},
		{
			name:         "service returns conflict",
			collectionID: collectionID,
			venueID:      venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("AddVenue", mock.Anything, collectionID, customerID, venueID).
					Return(apperror.ErrVenueAlreadyInCollection).
					Once()
			},
			wantCode: http.StatusConflict,
			wantBody: response.ErrorResponse{Message: apperror.ErrVenueAlreadyInCollection.Error()},
		},
		{
			name:         "service returns not found",
			collectionID: collectionID,
			venueID:      venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				err := apperror.CollectionNotFoundID(collectionID)
				s.On("AddVenue", mock.Anything, collectionID, customerID, venueID).
					Return(err).
					Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: response.ErrorResponse{Message: apperror.CollectionNotFoundID(collectionID).Error()},
		},
		{
			name:         "service returns unknown error",
			collectionID: collectionID,
			venueID:      venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("AddVenue", mock.Anything, collectionID, customerID, venueID).
					Return(errors.New("db is down")).
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

			r.POST("/collections/:collectionId/venues/:venueId", withCustomerID(tt.customerID, h.addVenue))

			w := performRequest(
				r,
				http.MethodPost,
				"/collections/"+tt.collectionID+"/venues/"+tt.venueID,
			)

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
