package business

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/business/error/code"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockParser struct{}

func (m *mockParser) ParseAccessToken(token string) (string, string, error) {
	return "1", "business", nil
}

type mockBusinessService struct {
	getFunc    func(ctx context.Context, id int) (*entity.OrgUnit, error)
	listFunc   func(ctx context.Context, brandId, page, limit int) ([]entity.OrgUnit, error)
	updateFunc func(ctx context.Context, postID int64, userID int64, content string) (*entity.Post, error)
	deleteFunc func(ctx context.Context, postID int64, userID int64) error
}

func (m *mockBusinessService) Get(ctx context.Context, id int) (*entity.OrgUnit, error) {
	return m.getFunc(ctx, id)
}

func (m *mockBusinessService) List(ctx context.Context, brandId, page, limit int) ([]entity.OrgUnit, error) {
	return m.listFunc(ctx, brandId, page, limit)
}

func (m *mockBusinessService) UpdatePost(ctx context.Context, postID int64, userID int64, content string) (*entity.Post, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, postID, userID, content)
	}
	return nil, nil
}

func (m *mockBusinessService) DeletePost(ctx context.Context, postID int64, userID int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, postID, userID)
	}
	return nil
}

func ptr[T any](v T) *T { return &v }

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		last := c.Errors.Last()
		if last == nil {
			return
		}

		var appErr *apperror.Error
		if errors.As(last.Err, &appErr) {
			switch appErr.Code {
			case code.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": appErr.Error()})
				return
			case code.BadRequest:
				c.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func setupRouter(svc *mockBusinessService) *gin.Engine {
	r := gin.New()
	r.Use(errorMiddleware())
	g := r.Group("/business")
	parser := &mockParser{}
	RegisterHandlers(g, svc, parser)
	return r
}

func TestIntegration_Get_OK(t *testing.T) {
	svc := &mockBusinessService{
		getFunc: func(_ context.Context, id int) (*entity.OrgUnit, error) {
			if id == 1 {
				return &entity.OrgUnit{
					Id: 1, Name: "Location 1", Avatar: ptr("a.png"),
					ProfileType: "VENUE", ParentId: ptr(2),
				}, nil
			}
			return &entity.OrgUnit{Id: 2, Name: "Brand", Avatar: ptr("b.png"), ProfileType: "BRAND"}, nil
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/business/1", nil)
	setupRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body getResponse
	json.Unmarshal(w.Body.Bytes(), &body)
	if body.Name != "Location 1" {
		t.Errorf("unexpected name: %v", body.Name)
	}
	if body.Brand == nil {
		t.Fatal("expected brand to be present")
	}
	if body.Brand.Name != "Brand" {
		t.Errorf("unexpected brand name: %v", body.Brand.Name)
	}
}

func TestIntegration_Get_NotFound(t *testing.T) {
	svc := &mockBusinessService{
		getFunc: func(_ context.Context, id int) (*entity.OrgUnit, error) {
			return nil, apperror.OrgUnitNotFoundID(id)
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/business/999", nil)
	setupRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIntegration_List_OK(t *testing.T) {
	svc := &mockBusinessService{
		listFunc: func(_ context.Context, brandId, page, limit int) ([]entity.OrgUnit, error) {
			if brandId != 1 {
				t.Errorf("expected brandId 1, got %d", brandId)
			}
			return []entity.OrgUnit{
				{Id: 10, Name: "Location A"},
				{Id: 11, Name: "Location B"},
			}, nil
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/business/1/locations?page=1&limit=10", nil)
	setupRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body listResponse
	json.Unmarshal(w.Body.Bytes(), &body)
	if len(body.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(body.Items))
	}
}

func TestIntegration_List_ServiceError(t *testing.T) {
	svc := &mockBusinessService{
		listFunc: func(_ context.Context, brandId, page, limit int) ([]entity.OrgUnit, error) {
			return nil, errors.New("db down")
		},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/business/1/locations?page=1&limit=10", nil)
	setupRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body.String())
	}
}
