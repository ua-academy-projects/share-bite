package business

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
	// "github.com/ua-academy-projects/share-bite/pkg/middleware"
)

type MockBusinessService struct {
	businessService
}

type dummyTokenParser struct{}

func (d dummyTokenParser) ParseAccessToken(token string) (string, string, error) {
	return "", "", nil
}

func (m *MockBusinessService) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int) (pagination.Result[entity.BoxWithDistance], error) {
	fakeBox := entity.BoxWithDistance{
		Box: entity.Box{
			Id:      123,
			VenueID: 1,
		},
		Distance: 2.5,
	}
	return pagination.Result[entity.BoxWithDistance]{
		Items: []entity.BoxWithDistance{fakeBox},
		Total: 1,
	}, nil
}

func TestListNearbyBoxes_InvalidCoordinates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockBusinessService{}
	parser := dummyTokenParser{}

	RegisterHandlers(router.Group("/"), mockService, parser)

	req, _ := http.NewRequest(http.MethodGet, "/nearby-boxes?lat=999&lon=30.5&limit=10", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got: %d", w.Code)
	}
}

func TestListNearbyBoxes_Positive_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockBusinessService{}
	parser := dummyTokenParser{}

	RegisterHandlers(router.Group("/"), mockService, parser)

	req, _ := http.NewRequest(http.MethodGet, "/nearby-boxes?lat=50.45&lon=30.52&limit=10&skip=0", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected: 200, got: %d. Body: %s", w.Code, w.Body.String())
	}

	responseBody := w.Body.String()

	if !strings.Contains(responseBody, `"total":1`) {
		t.Errorf("Expected total: 1 in response, got: %s", responseBody)
	}

	if !strings.Contains(responseBody, `"distance":2.5`) {
		t.Errorf("Expected distance: 2.5 in JSON, got: %s", responseBody)
	}
}
