package post

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	internalmiddleware "github.com/ua-academy-projects/share-bite/internal/middleware"
)

func TestPostHandler_Create(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{
		createFn: func(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
			return entity.Post{
				ID:         "post-1",
				CustomerID: in.CustomerID,
				VenueID:    in.VenueID,
				Text:       in.Text,
				Rating:     in.Rating,
				Status:     entity.PostStatusPublished,
				CreatedAt:  time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC),
			}, nil
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

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.WriteField("venue_id", "123"))
	require.NoError(t, writer.WriteField("text", "nice food"))
	require.NoError(t, writer.WriteField("rating", "5"))
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/posts/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusCreated, res.Code)
	assert.Equal(t, "user-1", customerSvc.lastUserID)
	assert.Equal(t, "customer-1", postSvc.lastCreateInput.CustomerID)
	assert.Equal(t, int64(123), postSvc.lastCreateInput.VenueID)
	assert.Equal(t, "nice food", postSvc.lastCreateInput.Text)
	assert.Equal(t, int16(5), postSvc.lastCreateInput.Rating)

	var got createResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
	assert.Equal(t, "post-1", got.Post.ID)
}

func TestPostHandler_Create_UnauthorizedWithoutHeader(t *testing.T) {
	t.Parallel()

	authMiddleware := internalmiddleware.Auth(tokenParserMock{})
	router := testRouter(&postServiceMock{}, &customerServiceMock{}, authMiddleware)

	body := `{"venue_id":123,"text":"nice food","rating":5}`
	req := httptest.NewRequest(http.MethodPost, "/posts/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Contains(t, res.Body.String(), "empty auth header")
}

func TestPostHandler_Create_InvalidPayload(t *testing.T) {
	t.Parallel()

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

	router := testRouter(&postServiceMock{}, customerSvc, authMiddleware)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.WriteField("venue_id", "123"))
	require.NoError(t, writer.WriteField("rating", "5"))
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/posts/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)

	assertResponseMessageNotEmpty(t, res)
}

func TestPostHandler_Create_InvalidContentType(t *testing.T) {
	t.Parallel()

	authMiddleware := internalmiddleware.Auth(tokenParserMock{
		parseAccessTokenFn: func(token string) (string, string, error) {
			if token != "valid-token" {
				return "", "", context.Canceled
			}

			return "user-1", "customer", nil
		},
	})

	router := testRouter(&postServiceMock{}, &customerServiceMock{}, authMiddleware)

	req := httptest.NewRequest(http.MethodPost, "/posts/", strings.NewReader(`{"venue_id":123,"text":"nice food","rating":5}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	assertResponseMessageContains(t, res, "multipart/form-data")
}
