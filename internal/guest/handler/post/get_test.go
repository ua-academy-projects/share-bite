package post

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func TestPostHandler_Get(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{
				ID:     postID,
				Text:   "hello",
				Rating: 4,
				Status: entity.PostStatusPublished,
				Images: []entity.PostImage{{
					ObjectKey: "posts/42/cover.jpg",
				}},
			}, nil
		},
	}

	router := testRouter(postSvc, &customerServiceMock{}, func(c *gin.Context) { c.Next() })

	req := httptest.NewRequest(http.MethodGet, "/posts/42", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "42", postSvc.lastGetID)

	var got getResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
	assert.Equal(t, "42", got.Post.ID)
	require.Len(t, got.Post.Images, 1)
	assert.Equal(t, "https://cdn.example/posts/42/cover.jpg", got.Post.Images[0])
}

func TestPostHandler_Get_ServiceError(t *testing.T) {
	t.Parallel()

	postSvc := &postServiceMock{
		getFn: func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
			return entity.Post{}, apperror.PostNotFoundID(postID)
		},
	}

	router := testRouter(postSvc, &customerServiceMock{}, func(c *gin.Context) { c.Next() })

	req := httptest.NewRequest(http.MethodGet, "/posts/42", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusNotFound, res.Code)

	assertResponseMessageContains(t, res, "post with id")
}
