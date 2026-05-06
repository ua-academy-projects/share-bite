package customer

import (
	"errors"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

func TestGetByUserName(t *testing.T) {
	t.Parallel()

	var (
		userName = "sharebite04"

		userID     = gofakeit.UUID()
		customerID = gofakeit.UUID()
		firstName  = gofakeit.FirstName()
		lastName   = gofakeit.LastName()
		bio        = gofakeit.Bio()

		baseURL         = "https://cdn.com/"
		avatarObjectKey = gofakeit.Word()
	)

	tests := []struct {
		name string

		userName string
		mockFn   func(s *mockCustomerService, st *mockObjectStorage)

		wantBody any
		wantCode int
	}{
		{
			name:     "success",
			userName: userName,
			mockFn: func(s *mockCustomerService, st *mockObjectStorage) {
				s.On("GetByUserName", mock.Anything, userName).
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

				st.On("BuildURL", avatarObjectKey).
					Return(baseURL + avatarObjectKey).
					Once()
			},
			wantBody: getByUserNameResponse{
				Customer: customerResponse{
					ID:        customerID,
					UserID:    userID,
					UserName:  userName,
					FirstName: firstName,
					LastName:  lastName,
					Bio:       &bio,
					AvatarURL: strPtr(baseURL + avatarObjectKey),
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "validation error - username too short",
			userName: "ab",
			mockFn:   func(s *mockCustomerService, st *mockObjectStorage) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "username", Message: "This field must be at least 3 characters long"},
				},
			},
		},
		{
			name:     "validation error - username too long",
			userName: "thisusernameiswaytoolongtoobevalid04",
			mockFn:   func(s *mockCustomerService, st *mockObjectStorage) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "username", Message: "This field must be at most 30 characters long"},
				},
			},
		},
		{
			name:     "validation error - invalid characters",
			userName: "sharebite!",
			mockFn:   func(s *mockCustomerService, st *mockObjectStorage) {},
			wantCode: http.StatusBadRequest,
			wantBody: response.ErrorResponse{
				Message: validationMsg,
				Details: []response.ErrorDetail{
					{Field: "username", Message: "This field can only contain letters and numbers"},
				},
			},
		},
		{
			name:     "error - customer not found",
			userName: "sharebitefake",
			mockFn: func(s *mockCustomerService, st *mockObjectStorage) {
				s.On("GetByUserName", mock.Anything, "sharebitefake").
					Return(entity.Customer{}, apperror.CustomerNotFoundUserName("sharebitefake")).
					Once()
			},
			wantCode: http.StatusNotFound,
			wantBody: response.ErrorResponse{Message: apperror.CustomerNotFoundUserName("sharebitefake").Error()},
		},
		{
			name:     "service unknown error",
			userName: userName,
			mockFn: func(s *mockCustomerService, st *mockObjectStorage) {
				s.On("GetByUserName", mock.Anything, userName).
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
			storage := new(mockObjectStorage)
			h := &handler{service: customerService, storage: storage}
			tt.mockFn(customerService, storage)

			r := newTestRouter()
			r.GET("/customers/:username", h.getByUserName)

			w := performRequest(r, http.MethodGet, "/customers/"+tt.userName)

			require.Equal(t, tt.wantCode, w.Code)
			if tt.wantBody != nil {
				assertJSONBody(t, tt.wantBody, w.Body.String())
			}

			customerService.AssertExpectations(t)
			storage.AssertExpectations(t)
		})
	}
}
