package customer

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/mock"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		userID     = gofakeit.UUID()
		customerID = gofakeit.UUID()

		userName  = "sharebite04"
		firstName = gofakeit.FirstName()
		lastName  = gofakeit.LastName()
		bio       = gofakeit.Bio()
	)

	tests := []struct {
		name string

		userID any
		body   any

		mockFn func(s *mockCustomerService)

		wantBody any
		wantCode int
	}{
		{
			name:   "success",
			userID: userID,
			body: createRequest{
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
				Bio:       strPtr(bio),
			},
			mockFn: func(s *mockCustomerService) {
				s.On("Create", mock.Anything, entity.CreateCustomer{
					UserID:    userID,
					UserName:  userName,
					FirstName: firstName,
					LastName:  lastName,
					Bio:       strPtr(bio),
				}).
					Return(customerID, nil).
					Once()
			},
			wantBody: createResponse{CustomerID: customerID},
			wantCode: http.StatusCreated,
		},
		{
			name:     "invalid json",
			body:     "{broken-json",
			userID:   userID,
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{Message: apperror.ErrInvalidJSON.Error()},
		},
		{
			name: "binding validation error",
			body: createRequest{
				UserName:  "",
				FirstName: "",
				LastName:  "",
			},
			userID:   userID,
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "userName", Message: "This field is required"},
					{Field: "firstName", Message: "This field is required"},
					{Field: "lastName", Message: "This field is required"},
				},
			},
		},
		{
			name: "custom mapper validation error - spaces only",
			body: createRequest{
				UserName:  userName,
				FirstName: "   ",
				LastName:  lastName,
			},
			userID:   userID,
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "firstName", Message: "This field must be at least 2 characters long"},
				},
			},
		},
		{
			name: "error - username taken",
			body: createRequest{
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
			},
			userID: userID,
			mockFn: func(s *mockCustomerService) {
				s.On("Create", mock.Anything, entity.CreateCustomer{
					UserID:    userID,
					UserName:  userName,
					FirstName: firstName,
					LastName:  lastName,
				}).
					Return("", apperror.CustomerUserNameTaken(userName)).
					Once()
			},
			wantCode: http.StatusConflict,
			wantBody: response.ErrorResponse{Message: apperror.CustomerUserNameTaken(userName).Error()},
		},
		{
			name: "unauthorized - no user id in context",
			body: createRequest{
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
			},
			userID:   nil,
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusInternalServerError, // we expect to have token in context after auth middleware step!
			wantBody: response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name: "service unknown error",
			body: createRequest{
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
				Bio:       strPtr(bio),
			},
			userID: userID,
			mockFn: func(s *mockCustomerService) {
				s.On("Create", mock.Anything, entity.CreateCustomer{
					UserID:    userID,
					UserName:  userName,
					FirstName: firstName,
					LastName:  lastName,
					Bio:       strPtr(bio),
				}).
					Return("", errors.New("service unexpected error")).
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
			r.POST("/customers", withUserID(tt.userID, h.create))

			var w *httptest.ResponseRecorder
			if s, ok := tt.body.(string); ok {
				w = performRawJSONRequest(t, r, http.MethodPost, "/customers", s)
			} else {
				w = performJSONRequest(t, r, http.MethodPost, "/customers", tt.body)
			}

			require.Equal(t, tt.wantCode, w.Code)
			assertJSONBody(t, tt.wantBody, w.Body.String())

			customerService.AssertExpectations(t)
		})
	}
}
