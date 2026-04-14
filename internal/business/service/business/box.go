package business

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

const (
	kilometerIndex = 1.60934
)

func (s *service) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int) (pagination.Result[entity.BoxWithDistance], error) {
	result, err := s.businessRepo.ListNearbyBoxes(ctx, offset, limit, lat, lon, categoryID)
	if err != nil {
		return pagination.Result[entity.BoxWithDistance]{}, fmt.Errorf("list nearby boxes: %w", err)
	}
	for i := range result.Items{
		result.Items[i].Distance = result.Items[i].Distance * kilometerIndex
	}
	return result, nil
}
