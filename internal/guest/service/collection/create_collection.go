package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) CreateCollection(ctx context.Context, in entity.CreateCollectionInput) (entity.Collection, error) {
	collection, err := s.collectionRepo.CreateCollection(ctx, in)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("create collection in repository: %w", err)
	}

	return collection, nil
}
