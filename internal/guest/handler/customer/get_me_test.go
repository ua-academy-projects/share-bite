package customer

import (
	"errors"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestGetMe(t *testing.T) {
	t.Parallel()

	var (
		userID     = gofakeit.UUID()
		customerID = gofakeit.UUID()

		userName        = "sharebite04"
		firstName       = gofakeit.FirstName()
		lastName        = gofakeit.LastName()
		avatarObjectKey = gofakeit.Word()
		bio             = gofakeit.Bio()
	)

	tests := []struct {
		name string

		userID any
		mockFn func(s *mockCustomerService)

		wantBody any
		wantCode int
	}{
		{
			name:   "success",
			userID: userID,
			mockFn: func(s *mockCustomerService) {
				s.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{
						ID:              customerID,
						UserID:          userID,
						UserName:        userName,
						FirstName:       firstName,
						LastName:        lastName,
						AvatarObjectKey: &avatarObjectKey,
						Bio:             &bio,
					}, nil).
					Once()
			},
			wantBody: getByUserNameResponse{
				Customer: customerToResponse(entity.Customer{
					ID:              customerID,
					UserID:          userID,
					UserName:        userName,
					FirstName:       firstName,
					LastName:        lastName,
					AvatarObjectKey: &avatarObjectKey,
					Bio:             &bio,
				}),
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "error - customer not found",
			userID: userID,
			mockFn: func(s *mockCustomerService) {
				s.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{}, apperror.CustomerNotFoundUserID(userID)).
					Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: response.ErrorResponse{Message: apperror.CustomerNotFoundUserID(userID).Error()},
		},
		{
			name:     "unauthorized - no user id in context",
			userID:   nil,
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusInternalServerError, // we expect to have token in context after auth middleware step!
			wantBody: response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:   "service unknown error",
			userID: userID,
			mockFn: func(s *mockCustomerService) {
				s.On("GetByUserID", mock.Anything, userID).
					Return(entity.Customer{}, errors.New("database connection lost")).
					Once()
			},
			wantCode: http.StatusInternalServerError,
			wantBody: response.ErrorResponse{Message: internalErrMsg},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			customerService := new(mockCustomerService)
			h := &handler{service: customerService}
			tt.mockFn(customerService)

			r := newTestRouter()
			r.GET("/customers", withUserID(tt.userID, h.getMe))
			w := performRequest(r, http.MethodGet, "/customers")

			require.Equal(t, tt.wantCode, w.Code)
			if tt.wantBody != nil {
				assertJSONBody(t, tt.wantBody, w.Body.String())
			}

			customerService.AssertExpectations(t)
		})
	}
}
