package customer

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	var (
		userID     = gofakeit.UUID()
		customerID = gofakeit.UUID()

		userName  = "sharebite04"
		firstName = gofakeit.FirstName()
		lastName  = gofakeit.LastName()
		bio       = gofakeit.Bio()
		avatarKey = gofakeit.Word()

		expectedCustomer = entity.Customer{
			ID:              customerID,
			UserID:          userID,
			UserName:        userName,
			FirstName:       firstName,
			LastName:        lastName,
			Bio:             &bio,
			AvatarObjectKey: &avatarKey,
		}
	)

	tests := []struct {
		name string

		input  any
		userID any

		mockFn func(s *mockCustomerService)

		wantBody any
		wantCode int
	}{
		{
			name:   "success",
			userID: userID,
			input: updateRequest{
				UserName:        strPtr(userName),
				FirstName:       strPtr(firstName),
				LastName:        strPtr(lastName),
				Bio:             strPtr(bio),
				AvatarObjectKey: strPtr(avatarKey),
			},
			mockFn: func(s *mockCustomerService) {
				s.On("Update", mock.Anything, entity.UpdateCustomer{
					UserID:          userID,
					UserName:        strPtr(userName),
					FirstName:       strPtr(firstName),
					LastName:        strPtr(lastName),
					Bio:             strPtr(bio),
					AvatarObjectKey: strPtr(avatarKey),
				}).
					Return(expectedCustomer, nil).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: updateResponse{Customer: customerToResponse(expectedCustomer)},
		},
		{
			name:   "success with trimming and lowercasing",
			userID: userID,
			input: updateRequest{
				UserName:  strPtr("ShareBite04"),
				FirstName: strPtr("  " + firstName + "  "),
				LastName:  strPtr("  " + lastName + "  "),
			},
			mockFn: func(s *mockCustomerService) {
				s.On("Update", mock.Anything, entity.UpdateCustomer{
					UserID:          userID,
					UserName:        strPtr("sharebite04"),
					FirstName:       strPtr(firstName),
					LastName:        strPtr(lastName),
					Bio:             nil,
					AvatarObjectKey: nil,
				}).
					Return(expectedCustomer, nil).
					Once()
			},
			wantCode: http.StatusOK,
			wantBody: updateResponse{Customer: customerToResponse(expectedCustomer)},
		},
		{
			name:     "invalid json",
			input:    "{broken-json",
			userID:   userID,
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{Message: apperror.ErrInvalidJSON.Error()},
		},
		{
			name:   "binding validation error - username too short",
			userID: userID,
			input: updateRequest{
				UserName: strPtr("ab"),
			},
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "userName", Message: "This field must be at least 3 characters long"},
				},
			},
		},
		{
			name:   "custom mapper validation error - spaces only",
			userID: userID,
			input: updateRequest{
				FirstName: strPtr("   "),
			},
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
			name:   "error - empty update",
			userID: userID,
			input:  updateRequest{},
			mockFn: func(s *mockCustomerService) {
				s.On("Update", mock.Anything, mock.AnythingOfType("entity.UpdateCustomer")).
					Return(entity.Customer{}, apperror.ErrEmptyUpdate).
					Once()
			},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{Message: apperror.ErrEmptyUpdate.Error()},
		},
		{
			name:   "error - username taken",
			userID: userID,
			input: updateRequest{
				UserName: strPtr(userName),
			},
			mockFn: func(s *mockCustomerService) {
				s.On("Update", mock.Anything, mock.AnythingOfType("entity.UpdateCustomer")).
					Return(entity.Customer{}, apperror.CustomerUserNameTaken(userName)).
					Once()
			},
			wantCode: http.StatusConflict,
			wantBody: response.ErrorResponse{Message: apperror.CustomerUserNameTaken(userName).Error()},
		},
		{
			name:     "unauthorized - no user id in context",
			userID:   nil,
			input:    updateRequest{FirstName: strPtr(firstName)},
			mockFn:   func(s *mockCustomerService) {},
			wantCode: http.StatusInternalServerError,
			wantBody: response.ErrorResponse{Message: internalErrMsg},
		},
		{
			name:   "service unknown error",
			userID: userID,
			input: updateRequest{
				FirstName: strPtr(firstName),
			},
			mockFn: func(s *mockCustomerService) {
				s.On("Update", mock.Anything, mock.AnythingOfType("entity.UpdateCustomer")).
					Return(entity.Customer{}, errors.New("unexpected db error")).
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
			r.PATCH("/customers", withUserID(tt.userID, h.update))

			var w *httptest.ResponseRecorder
			if s, ok := tt.input.(string); ok {
				w = performRawJSONRequest(t, r, http.MethodPatch, "/customers", s)
			} else {
				w = performJSONRequest(t, r, http.MethodPatch, "/customers", tt.input)
			}

			require.Equal(t, tt.wantCode, w.Code)

			assertJSONBody(t, tt.wantBody, w.Body.String())

			customerService.AssertExpectations(t)
		})
	}
}
