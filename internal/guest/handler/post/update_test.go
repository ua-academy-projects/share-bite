package post

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	internalmiddleware "github.com/ua-academy-projects/share-bite/internal/middleware"
)

func TestPostHandler_Update(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{
		updateFn: func(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
			text := "updated text"
			venueID := int64(456)
			status := entity.PostStatusPublished
			if in.VenueID != nil {
				venueID = *in.VenueID
			}
			if in.Status != nil {
				status = *in.Status
			}
			return entity.Post{
				ID:         in.ID,
				CustomerID: "customer-1",
				VenueID:    venueID,
				Text:       text,
				Rating:     4,
				Status:     status,
				CreatedAt:  time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2026, 4, 12, 11, 0, 0, 0, time.UTC),
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
	require.NoError(t, writer.WriteField("venue_id", "456"))
	require.NoError(t, writer.WriteField("text", "updated text"))
	require.NoError(t, writer.WriteField("rating", "4"))
	require.NoError(t, writer.WriteField("status", "archived"))
	imagePart, err := writer.CreateFormFile("images", "cover.png")
	require.NoError(t, err)
	_, err = imagePart.Write([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'})
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPatch, "/posts/42", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "42", postSvc.lastUpdateInput.ID)
	assert.Equal(t, "user-1", customerSvc.lastUserID)
	assert.Equal(t, "customer-1", postSvc.lastUpdateInput.CustomerID)
	if assert.NotNil(t, postSvc.lastUpdateInput.VenueID) {
		assert.Equal(t, int64(456), *postSvc.lastUpdateInput.VenueID)
	}
	if assert.NotNil(t, postSvc.lastUpdateInput.Text) {
		assert.Equal(t, "updated text", *postSvc.lastUpdateInput.Text)
	}
	if assert.NotNil(t, postSvc.lastUpdateInput.Rating) {
		assert.Equal(t, int16(4), *postSvc.lastUpdateInput.Rating)
	}
	if assert.NotNil(t, postSvc.lastUpdateInput.Status) {
		assert.Equal(t, entity.PostStatusArchived, *postSvc.lastUpdateInput.Status)
	}
	assert.True(t, postSvc.lastUpdateInput.RewriteImages)
	require.Len(t, postSvc.lastUpdateInput.Images, 1)
}

func TestPostHandler_Update_UnauthorizedWithoutHeader(t *testing.T) {
	t.Parallel()

	authMiddleware := internalmiddleware.Auth(tokenParserMock{})
	router := testRouter(&postServiceMock{}, &customerServiceMock{}, authMiddleware)

	body := `{"text":"updated text"}`
	req := httptest.NewRequest(http.MethodPatch, "/posts/42", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnauthorized, res.Code)
	assert.Contains(t, res.Body.String(), "empty auth header")
}

func TestPostHandler_Update_RejectsDeletedStatus(t *testing.T) {
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

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	require.NoError(t, writer.WriteField("status", "deleted"))
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPatch, "/posts/42", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)

	assertResponseMessageContains(t, res, "deleted")

	assert.Empty(t, postSvc.lastUpdateInput.ID)
}

func TestPostHandler_Update_InvalidContentType(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPatch, "/posts/42", strings.NewReader(`{"text":"updated"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)

	assertResponseMessageContains(t, res, "multipart/form-data")
}
