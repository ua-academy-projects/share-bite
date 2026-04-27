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

func TestRemoveVenue(t *testing.T) {
	t.Parallel()

	var (
		collectionID        = gofakeit.UUID()
		invalidCollectionID = "bad-uuid"
		customerID          = gofakeit.UUID()
		venueID             = int64(gofakeit.Number(1, 1_000_000))
		venueIDStr          = strconv.FormatInt(venueID, 10)
	)

	tests := []struct {
		name string

		collectionID string
		venueIDPath  string
		customerID   any
		mockFn       func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("RemoveVenue", mock.Anything, collectionID, customerID, venueID).
					Return(nil).
					Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:         "validation error invalid collection id",
			collectionID: invalidCollectionID,
			venueIDPath:  venueIDStr,
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
			name:         "validation error invalid venue id",
			collectionID: collectionID,
			venueIDPath:  "0",
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
			name:         "validation error negative venue id",
			collectionID: collectionID,
			venueIDPath:  "-4",
			customerID:   customerID,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "venueId", Message: "This field must be greater than or equal to 1"},
				},
			},
		},
		{
			name:         "missing customer id in ctx",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   nil,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service venue not found in collection",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("RemoveVenue", mock.Anything, collectionID, customerID, venueID).
					Return(apperror.VenueNotFoundInCollection(venueID)).
					Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: response.ErrorResponse{Message: apperror.VenueNotFoundInCollection(venueID).Error()},
		},
		{
			name:         "service unknown error",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("RemoveVenue", mock.Anything, collectionID, customerID, venueID).
					Return(errors.New("boom")).
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

			r.DELETE("/collections/:collectionId/venues/:venueId", withCustomerID(tt.customerID, h.removeVenue))

			w := performRequest(r, http.MethodDelete, "/collections/"+tt.collectionID+"/venues/"+tt.venueIDPath)

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
