package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) GetCollection(ctx context.Context, collectionID string, customerID *string) (entity.Collection, error) {
	collection, err := s.collectionRepo.GetCollection(ctx, collectionID)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("get collection from repository: %w", err)
	}

	// if its not my collection and not public one -> not found
	if !canAccessCollection(collection, customerID) {
		return entity.Collection{}, apperror.CollectionNotFoundID(collectionID)
	}

	return collection, nil
}
