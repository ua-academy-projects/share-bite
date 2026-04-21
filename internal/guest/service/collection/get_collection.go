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

	if !collection.IsPublic {
		if customerID == nil {
			return entity.Collection{}, apperror.CollectionNotFoundID(collectionID)
		}

		if err := s.requireCollaborator(ctx, collectionID, *customerID, collection.CustomerID); err != nil {
			return entity.Collection{}, err
		}
	}

	return collection, nil
}
