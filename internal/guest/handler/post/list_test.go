package post

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func TestPostHandler_List(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(10 * time.Minute)

	postSvc := &postServiceMock{
		listFn: func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
			return dto.ListPostsOutput{
				Posts: []entity.Post{{
					ID:         "1",
					CustomerID: "c-1",
					VenueID:    101,
					Text:       "great place",
					Rating:     5,
					Status:     entity.PostStatusPublished,
					CreatedAt:  createdAt,
					UpdatedAt:  updatedAt,
					Images: []entity.PostImage{{
						ObjectKey: "posts/1/main.jpg",
					}, {
						ObjectKey: "posts/1/second.jpg",
					}},
				}},
				Total: 1,
			}, nil
		},
	}

	router := testRouter(postSvc, &customerServiceMock{}, func(c *gin.Context) { c.Next() })

	req := httptest.NewRequest(http.MethodGet, "/posts/?limit=1&offset=0", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, 1, postSvc.lastListInput.Limit)
	assert.Equal(t, 0, postSvc.lastListInput.Offset)

	var got listResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&got))
	require.Len(t, got.Posts, 1)
	assert.Equal(t, 1, got.Total)
	assert.Equal(t, "1", got.Posts[0].ID)
	require.Len(t, got.Posts[0].Images, 2)
	assert.Equal(t, "https://cdn.example/posts/1/main.jpg", got.Posts[0].Images[0])
	assert.Equal(t, "https://cdn.example/posts/1/second.jpg", got.Posts[0].Images[1])
}

func TestPostHandler_List_InvalidQuery(t *testing.T) {
	t.Parallel()

	router := testRouter(&postServiceMock{}, &customerServiceMock{}, func(c *gin.Context) { c.Next() })

	req := httptest.NewRequest(http.MethodGet, "/posts/?limit=101", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)

	assertResponseMessageNotEmpty(t, res)
}
