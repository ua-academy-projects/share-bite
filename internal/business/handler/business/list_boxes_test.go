package business

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

type MockBusinessService struct {
	businessService
}

type dummyTokenParser struct{}

func (d dummyTokenParser) ParseAccessToken(token string) (string, string, error) {
	return "", "", nil
}

func (m *MockBusinessService) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int) (pagination.Result[entity.BoxWithDistance], error) {
	fakeBox1 := entity.BoxWithDistance{
		Box: entity.Box{
			Id:      123,
			VenueId: 1,
		},
		AvailabilityCount: 6,
		Distance: 2.5,
	}
	fakeBox2 := entity.BoxWithDistance{
		Box: entity.Box{
			Id:321,
			VenueId: 1,
		},
		AvailabilityCount: 8,
		Distance: 3,
	}
	return pagination.Result[entity.BoxWithDistance]{
		Items: []entity.BoxWithDistance{fakeBox1, fakeBox2},
		Total: 2,
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

	var response dto.ListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to Unmarshal resp: %v", err)
	}

	if response.Total != 2 {
		t.Errorf("Expected total: 2 in response, got: %d", response.Total)
	}

	if response.Items[0].AvailabilityStatus != "running_low" {
		t.Errorf("Expected status: running_low in response, got: %s", response.Items[0].AvailabilityStatus)
	}

	if response.Items[1].AvailabilityStatus != "available" {
		t.Errorf("Expected status: available in response, got: %s", response.Items[1].AvailabilityStatus)
	}
	
	if response.Items[0].Distance != 2.5 {
		t.Errorf("Expected distance: 2.5 in response, got: %f", response.Items[0].Distance)
	}
}