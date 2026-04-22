package collection

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestCreateCollection(t *testing.T) {
	t.Parallel()

	var (
		customerID = gofakeit.UUID()

		now  = time.Now().UTC()
		desc = gofakeit.ProductDescription()
	)

	tests := []struct {
		name string

		body          any
		customerIDVal any
		mockFn        func(s *mockCollectionService)

		wantCode int
		wantBody any
	}{
		{
			name:          "success",
			body:          createCollectionRequest{Name: "  My collection  ", Description: &desc, IsPublic: boolPtr(true)},
			customerIDVal: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("CreateCollection", mock.Anything, entity.CreateCollectionInput{
					CustomerID:  customerID,
					Name:        "My collection",
					Description: &desc,
					IsPublic:    true,
				}).Return(entity.Collection{
					ID:          "e7951b01-65d6-475b-b8e8-1765a67464af",
					CustomerID:  customerID,
					Name:        "My collection",
					Description: &desc,
					IsPublic:    true,
					CreatedAt:   now,
					UpdatedAt:   now,
				}, nil).Once()
			},
			wantCode: http.StatusCreated,
			wantBody: createCollectionResponse{
				Collection: collectionResponse{
					ID:          "e7951b01-65d6-475b-b8e8-1765a67464af",
					Name:        "My collection",
					Description: &desc,
					IsPublic:    true,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name:          "invalid json",
			body:          "{broken-json",
			customerIDVal: customerID,
			mockFn:        func(s *mockCollectionService) {},
			wantCode:      http.StatusBadRequest,
			wantBody:      response.ErrorResponse{Message: apperror.ErrInvalidJSON.Error()},
		},
		{
			name:          "validation error",
			body:          createCollectionRequest{Name: "", IsPublic: nil},
			customerIDVal: customerID,
			mockFn:        func(s *mockCollectionService) {},
			wantCode:      http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "name", Message: "This field is required"},
					{Field: "isPublic", Message: "This field is required"},
				},
			},
		},
		{
			name:          "validation error - name must be at least 1 char long",
			body:          createCollectionRequest{Name: " ", IsPublic: boolPtr(true)},
			customerIDVal: customerID,
			mockFn:        func(s *mockCollectionService) {},
			wantCode:      http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "name", Message: "This field must be at least 1 characters long"},
				},
			},
		},
		{
			name:          "missing customer id in ctx",
			body:          createCollectionRequest{Name: "A", IsPublic: boolPtr(true)},
			customerIDVal: nil,
			mockFn:        func(s *mockCollectionService) {},
			wantCode:      http.StatusInternalServerError,
			wantBody:      response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:          "service collection access denied error",
			body:          createCollectionRequest{Name: "A", IsPublic: boolPtr(true)},
			customerIDVal: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("CreateCollection", mock.Anything, mock.Anything).
					Return(entity.Collection{}, apperror.ErrCollectionAccessDenied).Once()
			},
			wantCode: http.StatusForbidden,
			wantBody: response.ErrorResponse{Message: apperror.ErrCollectionAccessDenied.Error()},
		},
		{
			name:          "service unknown error",
			body:          createCollectionRequest{Name: "A", IsPublic: boolPtr(true)},
			customerIDVal: customerID,
			mockFn: func(s *mockCollectionService) {
				s.On("CreateCollection", mock.Anything, mock.Anything).
					Return(entity.Collection{}, errors.New("boom")).Once()
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

			r.POST("/collections", withCustomerID(tt.customerIDVal, h.createCollection))

			var w *httptest.ResponseRecorder
			if s, ok := tt.body.(string); ok {
				w = performRawJSONRequest(t, r, http.MethodPost, "/collections", s)
			} else {
				w = performJSONRequest(t, r, http.MethodPost, "/collections", tt.body)
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
