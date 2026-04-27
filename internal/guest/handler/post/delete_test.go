package post

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	internalmiddleware "github.com/ua-academy-projects/share-bite/internal/middleware"
)

func TestPostHandler_Delete_ServiceError(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{
		deleteFn: func(ctx context.Context, postID, customerID string) error {
			return apperror.PostNotFoundID(postID)
		},
	}
	customerSvc := &customerServiceMock{
		getByUserIDFn: func(ctx context.Context, userID string) (entity.Customer, error) {
			return entity.Customer{ID: "customer-1", UserID: userID}, nil
		},
	}

	authMiddleware := internalmiddleware.Auth(tokenParserMock{
		parseAccessTokenFn: func(token string) (string, string, error) {
			if token != "valid-token" {
				return "", "", context.Canceled
			}

			return "user-1", "customer", nil
		},
	})

	router := testRouter(postSvc, customerSvc, authMiddleware)

	req := httptest.NewRequest(http.MethodDelete, "/posts/42", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)

	assertResponseMessageContains(t, res, "post with id")
}

func TestPostHandler_Delete(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{}
	customerSvc := &customerServiceMock{
		getByUserIDFn: func(ctx context.Context, userID string) (entity.Customer, error) {
			return entity.Customer{ID: "customer-1", UserID: userID}, nil
		},
	}

	authMiddleware := internalmiddleware.Auth(tokenParserMock{
		parseAccessTokenFn: func(token string) (string, string, error) {
			if token != "valid-token" {
				return "", "", context.Canceled
			}

			return "user-1", "customer", nil
		},
	})

	router := testRouter(postSvc, customerSvc, authMiddleware)

	req := httptest.NewRequest(http.MethodDelete, "/posts/42", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusNoContent, res.Code)
	assert.Equal(t, "42", postSvc.lastDeletePostID)
	assert.Equal(t, "customer-1", postSvc.lastDeleteCustomerID)
}

func TestPostHandler_Delete_UnauthorizedWithoutHeader(t *testing.T) {
	t.Parallel()

	authMiddleware := internalmiddleware.Auth(tokenParserMock{})
	router := testRouter(&postServiceMock{}, &customerServiceMock{}, authMiddleware)

	req := httptest.NewRequest(http.MethodDelete, "/posts/42", nil)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Contains(t, res.Body.String(), "empty auth header")
}
