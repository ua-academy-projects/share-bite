package customer

import (
	"context"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) GetByIDs(ctx context.Context, ids []string) ([]entity.Customer, error) {
	return s.customerRepo.GetByIDs(ctx, ids)
}
