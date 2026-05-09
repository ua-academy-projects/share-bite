package collection

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestReorderVenue(t *testing.T) {
	t.Parallel()

	var (
		collectionID        = gofakeit.UUID()
		invalidCollectionID = "bad-uuid"
		customerID          = gofakeit.UUID()

		venueID    = int64(gofakeit.Number(1, 1_000_000))
		venueIDStr = strconv.FormatInt(venueID, 10)

		prevVenueID = int64(gofakeit.Number(1, 1_000_000))
		nextVenueID = int64(gofakeit.Number(1, 1_000_000))
	)

	tests := []struct {
		name string

		collectionID string
		venueIDPath  string
		customerID   any
		body         any

		mockFn func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success with prevVenueId",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			body: reorderVenueRequest{
				PrevVenueID: &prevVenueID,
			},
			mockFn: func(s *mockCollectionService) {
				s.On("ReorderVenue", mock.Anything, entity.ReorderVenueInput{
					CollectionID: collectionID,
					VenueID:      venueID,
					CustomerID:   customerID,
					PrevVenueID:  &prevVenueID,
					NextVenueID:  nil,
				}).Return(nil).Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:         "success with nextVenueId",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			body: reorderVenueRequest{
				NextVenueID: &nextVenueID,
			},
			mockFn: func(s *mockCollectionService) {
				s.On("ReorderVenue", mock.Anything, entity.ReorderVenueInput{
					CollectionID: collectionID,
					VenueID:      venueID,
					CustomerID:   customerID,
					PrevVenueID:  nil,
					NextVenueID:  &nextVenueID,
				}).Return(nil).Once()
			},
			wantCode: http.StatusNoContent,
			wantBody: nil,
		},
		{
			name:         "validation error invalid collection id",
			collectionID: invalidCollectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			body: reorderVenueRequest{
				PrevVenueID: &prevVenueID,
			},
			mockFn:   func(s *mockCollectionService) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: "request validation failed",
				Details: []response.ErrorDetail{
					{Field: "collectionId", Message: "This field must be a valid UUID"},
				},
			},
		},
		{
			name:         "validation error invalid venue id in path",
			collectionID: collectionID,
			venueIDPath:  "0",
			customerID:   customerID,
			body: reorderVenueRequest{
				PrevVenueID: &prevVenueID,
			},
			mockFn:   func(s *mockCollectionService) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: "request validation failed",
				Details: []response.ErrorDetail{
					{Field: "venueId", Message: "This field is required"},
				},
			},
		},
		{
			name:         "validation error both neighbors missing",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			body:         reorderVenueRequest{},
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: "request validation failed",
				Details: []response.ErrorDetail{
					{Field: "prevVenueId", Message: "This field is required if NextVenueID is missing"},
					{Field: "nextVenueId", Message: "This field is required if PrevVenueID is missing"},
				},
			},
		},
		{
			name:         "missing customer id in ctx",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   nil,
			body: reorderVenueRequest{
				PrevVenueID: &prevVenueID,
			},
			mockFn:   func(s *mockCollectionService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: response.ErrorResponse{Message: "internal server error"},
		},
		{
			name:         "service forbidden",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			body: reorderVenueRequest{
				PrevVenueID: &prevVenueID,
			},
			mockFn: func(s *mockCollectionService) {
				s.On("ReorderVenue", mock.Anything, mock.Anything).
					Return(apperror.ErrCollectionAccessDenied).
					Once()
			},
			wantCode: http.StatusForbidden,
			wantBody: response.ErrorResponse{Message: apperror.ErrCollectionAccessDenied.Error()},
		},
		{
			name:         "service unknown error",
			collectionID: collectionID,
			venueIDPath:  venueIDStr,
			customerID:   customerID,
			body: reorderVenueRequest{
				PrevVenueID: &prevVenueID,
			},
			mockFn: func(s *mockCollectionService) {
				s.On("ReorderVenue", mock.Anything, mock.Anything).
					Return(errors.New("boom")).
					Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: response.ErrorResponse{Message: "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := new(mockCollectionService)
			tt.mockFn(svc)

			h := &handler{service: svc}
			r := newTestRouter()

			r.POST("/collections/:collectionId/venues/:venueId/reorder", withCustomerID(tt.customerID, h.reorderVenue))

			w := performJSONRequest(
				t,
				r,
				http.MethodPost,
				"/collections/"+tt.collectionID+"/venues/"+tt.venueIDPath+"/reorder",
				tt.body,
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
