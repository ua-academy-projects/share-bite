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
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestListMyCollections(t *testing.T) {
	t.Parallel()

	var (
		customerID   = gofakeit.UUID()
		collectionID = gofakeit.UUID()
		now          = time.Now().UTC()

		validToken = generatePageToken(now, collectionID)
	)

	tests := []struct {
		name string

		query      string
		customerID any
		mockFn     func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:       "success",
			query:      "",
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListCustomerCollections", mock.Anything, entity.ListCustomerCollectionsInput{
					CustomerID: customerID,
					CursorTime: time.Time{},
					CursorID:   "",
					Limit:      21,
				}).Return(entity.ListCustomerCollectionsOutput{
					Collections: []entity.Collection{
						{
							ID:         collectionID,
							CustomerID: customerID,
							Name:       "First",
							IsPublic:   true,
							CreatedAt:  now,
							UpdatedAt:  now,
						},
					},
					NextCursorTime: &now,
					NextCursorID:   &collectionID,
				}, nil).Once()
			},
			wantCode: http.StatusOK,
			wantBody: listMyCollectionsResponse{
				Collections: []collectionResponse{
					{
						ID:        collectionID,
						Name:      "First",
						IsPublic:  true,
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				// Очікуємо вже закодований рядок
				NextPageToken: validToken,
			},
		},
		{
			name:       "validation error",
			query:      "?pageSize=101",
			customerID: customerID,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "pageSize", Message: "This field must be less than or equal to 100"},
				},
			},
		},
		{
			name:       "missing customer id in ctx",
			query:      "",
			customerID: nil,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusInternalServerError,
			wantBody:   response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:       "invalid customer id type in ctx",
			query:      "",
			customerID: 123,
			mockFn:     func(s *mockCollectionService) {},
			wantCode:   http.StatusInternalServerError,
			wantBody:   response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:       "service unknown error",
			query:      "?pageSize=10&pageToken=" + validToken,
			customerID: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("ListCustomerCollections", mock.Anything, entity.ListCustomerCollectionsInput{
					CustomerID: customerID,
					CursorTime: now,
					CursorID:   collectionID,
					Limit:      11,
				}).Return(entity.ListCustomerCollectionsOutput{}, errors.New("db down")).Once()
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

			r.GET("/collections/me", withCustomerID(tt.customerID, h.listMyCollections))

			w := performRequest(r, http.MethodGet, "/collections/me"+tt.query)

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
