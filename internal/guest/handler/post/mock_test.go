package post

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/guest/error/code"
)

type postServiceMock struct {
	createFn func(ctx context.Context, in dto.CreatePostInput) (entity.Post, error)
	updateFn func(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	deleteFn func(ctx context.Context, postID, customerID string) error
	listFn   func(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error)
	getFn    func(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
	likeFn   func(ctx context.Context, postID string, customerID string) error
	unlikeFn func(ctx context.Context, postID string, customerID string) error

	lastCreateInput      dto.CreatePostInput
	lastUpdateInput      entity.UpdatePostInput
	lastDeletePostID     string
	lastDeleteCustomerID string
	lastListInput        dto.ListPostsInput
	lastGetID            string
}

func (m *postServiceMock) Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	m.lastCreateInput = in
	if m.createFn != nil {
		return m.createFn(ctx, in)
	}

	return entity.Post{}, nil
}

func (m *postServiceMock) Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
	m.lastUpdateInput = in
	if m.updateFn != nil {
		return m.updateFn(ctx, in)
	}

	return entity.Post{}, nil
}

func (m *postServiceMock) Delete(ctx context.Context, postID, customerID string) error {
	m.lastDeletePostID = postID
	m.lastDeleteCustomerID = customerID
	if m.deleteFn != nil {
		return m.deleteFn(ctx, postID, customerID)
	}

	return nil
}

func (m *postServiceMock) List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
	m.lastListInput = in
	if m.listFn != nil {
		return m.listFn(ctx, in)
	}

	return dto.ListPostsOutput{}, nil
}

func (m *postServiceMock) Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
	m.lastGetID = postID
	if m.getFn != nil {
		return m.getFn(ctx, postID, reqCustomerID)
	}

	return entity.Post{}, nil
}

func (m *postServiceMock) Like(ctx context.Context, postID string, customerID string) error {
	if m.likeFn != nil {
		return m.likeFn(ctx, postID, customerID)
	}
	return nil
}

func (m *postServiceMock) Unlike(ctx context.Context, postID string, customerID string) error {
	if m.unlikeFn != nil {
		return m.unlikeFn(ctx, postID, customerID)
	}
	return nil
}

type customerServiceMock struct {
	getByUserIDFn func(ctx context.Context, userID string) (entity.Customer, error)

	lastUserID string
}

func (m *customerServiceMock) GetByUserID(ctx context.Context, userID string) (entity.Customer, error) {
	m.lastUserID = userID
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(ctx, userID)
	}

	return entity.Customer{}, nil
}

type tokenParserMock struct {
	parseAccessTokenFn func(token string) (string, string, error)
}

func (m tokenParserMock) ParseAccessToken(token string) (string, string, error) {
	if m.parseAccessTokenFn != nil {
		return m.parseAccessTokenFn(token)
	}

	return "", "", nil
}

type objectStorageMock struct{}

func (objectStorageMock) Upload(context.Context, string, string, io.Reader) (string, error) {
	return "", nil
}

func (objectStorageMock) Delete(context.Context, string) error {
	return nil
}

func (objectStorageMock) BuildURL(key string) string {
	return "https://cdn.example/" + key
}

func testRouter(postSvc postService, customerSvc customerService, authMiddleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(testErrorMiddleware())
	RegisterHandlers(router.Group("/posts"), postSvc, customerSvc, authMiddleware, objectStorageMock{})

	return router
}

func testErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		respCode := http.StatusInternalServerError
		message := "internal server error"

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			message = appErr.Error()
			switch appErr.Code {
			case code.NotFound:
				respCode = http.StatusNotFound
			case code.InvalidJSON, code.InvalidRequest, code.BadRequest, code.EmptyUpdate:
				respCode = http.StatusBadRequest
			case code.UpstreamError:
				respCode = http.StatusBadGateway
			case code.AlreadyExists:
				respCode = http.StatusConflict
			case code.Forbidden:
				respCode = http.StatusForbidden
			default:
				respCode = http.StatusInternalServerError
			}
		}

		c.JSON(respCode, gin.H{"message": message})
	}
}

func assertResponseMessageContains(t *testing.T, res *httptest.ResponseRecorder, want string) {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	msg, ok := body["message"].(string)
	require.True(t, ok, "expected message field in response body, got: %+v", body)
	assert.Contains(t, strings.ToLower(msg), strings.ToLower(want))
}

func assertResponseMessageNotEmpty(t *testing.T, res *httptest.ResponseRecorder) {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	msg, ok := body["message"].(string)
	require.True(t, ok, "expected message field in response body, got: %+v", body)
	assert.NotEmpty(t, strings.TrimSpace(msg))
}
