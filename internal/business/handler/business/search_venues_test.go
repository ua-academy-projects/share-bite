package business

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

type searchVenuesServiceMock struct {
	businessService

	called   bool
	gotQuery string
	gotSkip  int
	gotLimit int
	gotTags  []string

	result pagination.Result[entity.OrgUnit]
	err    error
}

func (m *searchVenuesServiceMock) SearchVenues(
	ctx context.Context,
	query string,
	skip, limit int,
	tags []string,
) (pagination.Result[entity.OrgUnit], error) {
	m.called = true
	m.gotQuery = query
	m.gotSkip = skip
	m.gotLimit = limit
	m.gotTags = append([]string(nil), tags...)

	return m.result, m.err
}

func testBusinessErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		last := c.Errors.Last()
		if last == nil {
			return
		}

		var appErr *apperror.Error
		if errors.As(last.Err, &appErr) {
			switch appErr.Code {
			case code.BadRequest:
				c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
				return
			case code.Forbidden:
				c.JSON(http.StatusForbidden, gin.H{"error": appErr.Error()})
				return
			case code.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
				return
			case code.Unauthorized:
				c.JSON(http.StatusUnauthorized, gin.H{"error": appErr.Error()})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func setupSearchVenuesRouter(s businessService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(testBusinessErrorMiddleware())

	RegisterHandlers(r.Group("/"), s, dummyTokenParser{}, nil, nil)

	return r
}

func TestSearchVenues_EmptyFilters_ReturnsBadRequest(t *testing.T) {
	mock := &searchVenuesServiceMock{}
	router := setupSearchVenuesRouter(mock)

	req, _ := http.NewRequest(http.MethodGet, "/venues/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", w.Code, w.Body.String())
	}

	if mock.called {
		t.Fatalf("service should not be called when q and tags are empty")
	}
}

func TestSearchVenues_TagsAreNormalizedAndDeduplicated(t *testing.T) {
	mock := &searchVenuesServiceMock{
		result: pagination.Result[entity.OrgUnit]{
			Items: []entity.OrgUnit{},
			Total: 0,
		},
	}
	router := setupSearchVenuesRouter(mock)

	req, _ := http.NewRequest(http.MethodGet, "/venues/search?tags=%20Vegan%20,vegan,ROMANTIC", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	want := []string{"vegan", "romantic"}
	if !reflect.DeepEqual(mock.gotTags, want) {
		t.Fatalf("expected tags %v, got %v", want, mock.gotTags)
	}
}

func TestSearchVenues_DefaultPaginationApplied(t *testing.T) {
	mock := &searchVenuesServiceMock{
		result: pagination.Result[entity.OrgUnit]{
			Items: []entity.OrgUnit{},
			Total: 0,
		},
	}
	router := setupSearchVenuesRouter(mock)

	req, _ := http.NewRequest(http.MethodGet, "/venues/search?q=dinner", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	if mock.gotSkip != 0 {
		t.Fatalf("expected skip=0, got %d", mock.gotSkip)
	}
	if mock.gotLimit != 10 {
		t.Fatalf("expected limit=10, got %d", mock.gotLimit)
	}
	if mock.gotQuery != "dinner" {
		t.Fatalf("expected query=dinner, got %q", mock.gotQuery)
	}
}
