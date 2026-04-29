package collection

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestListVenues(t *testing.T) {
	t.Parallel()

	var (
		collectionID        = gofakeit.UUID()
		invalidCollectionID = "bad-uuid"
		customerID          = gofakeit.UUID()
		avatarURL           = gofakeit.URL()
		bannerURL           = gofakeit.URL()

		now  = time.Now().UTC()
		desc = gofakeit.ProductDescription()
	)

	tests := []struct {
		name string

		collectionID string
		customerID   any
		mockFn       func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:         "success with optional customer id",
			collectionID: collectionID,
			customerID:   customerID,
			mockFn: func(s *mockCollectionService) {
				cid := customerID
				s.On("ListVenues", mock.Anything, collectionID, &cid).
					Return([]entity.EnrichedVenueItem{
						{
							VenueItem: entity.Venue{
								ID:          10,
								Name:        "Venue 10",
								Description: strPtr(desc),
								AvatarURL:   strPtr(avatarURL),
								BannerURL:   strPtr(bannerURL),
							},
							SortOrder: 100,
							AddedAt:   now,
						},
					}, nil).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listVenuesResponse{
				Venues: []enrichedVenueItemResponse{
					{
						ID:          10,
						Name:        "Venue 10",
						Description: strPtr(desc),
						AvatarURL:   strPtr(avatarURL),
						BannerURL:   strPtr(bannerURL),
						SortOrder:   100,
						AddedAt:     now,
					},
				},
			},
		},
		{
			name:         "success without customer id",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				s.On("ListVenues", mock.Anything, collectionID, (*string)(nil)).
					Return([]entity.EnrichedVenueItem{}, nil).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: listVenuesResponse{Venues: []enrichedVenueItemResponse{}},
		},
		{
			name:         "validation error",
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
			name:         "invalid optional customer id type in ctx",
			collectionID: collectionID,
			customerID:   123,
			mockFn:       func(s *mockCollectionService) {},
			wantCode:     http.StatusInternalServerError,
			wantBody:     response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:         "service not found error",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				err := apperror.CollectionNotFoundID(collectionID)
				s.On("ListVenues", mock.Anything, collectionID, (*string)(nil)).
					Return([]entity.EnrichedVenueItem(nil), err).
					Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: response.ErrorResponse{Message: apperror.CollectionNotFoundID(collectionID).Error()},
		},
		{
			name:         "service unknown error",
			collectionID: collectionID,
			customerID:   nil,
			mockFn: func(s *mockCollectionService) {
				s.On("ListVenues", mock.Anything, collectionID, (*string)(nil)).
					Return([]entity.EnrichedVenueItem(nil), errors.New("boom")).
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

			r.GET("/collections/:collectionId/venues", withCustomerID(tt.customerID, h.listVenues))

			w := performRequest(r, http.MethodGet, "/collections/"+tt.collectionID+"/venues")

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
