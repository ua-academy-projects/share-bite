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

func TestPostHandler_Like(t *testing.T) {
	t.Parallel()

	liked := false
	postSvc := &postServiceMock{
		likeFn: func(ctx context.Context, postID string, customerID string) error {
			assert.Equal(t, "42", postID)
			assert.Equal(t, "customer-1", customerID)
			liked = true
			return nil
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

	req := httptest.NewRequest(http.MethodPost, "/posts/42/like", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusNoContent, res.Code)
	assert.True(t, liked)
}

func TestPostHandler_Like_ServiceError(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{
		likeFn: func(ctx context.Context, postID string, customerID string) error {
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

	req := httptest.NewRequest(http.MethodPost, "/posts/42/like", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
	assertResponseMessageContains(t, res, "post with id")
}

func TestPostHandler_Unlike(t *testing.T) {
	t.Parallel()

	unliked := false
	postSvc := &postServiceMock{
		unlikeFn: func(ctx context.Context, postID string, customerID string) error {
			assert.Equal(t, "42", postID)
			assert.Equal(t, "customer-1", customerID)
			unliked = true
			return nil
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

	req := httptest.NewRequest(http.MethodDelete, "/posts/42/like", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	assert.True(t, unliked)
}

func TestPostHandler_Unlike_ServiceError(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{
		unlikeFn: func(ctx context.Context, postID string, customerID string) error {
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

	req := httptest.NewRequest(http.MethodDelete, "/posts/42/like", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)
	assertResponseMessageContains(t, res, "post with id")
}
