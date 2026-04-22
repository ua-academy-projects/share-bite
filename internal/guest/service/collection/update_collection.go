package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error) {
	collection, err := s.collectionRepo.GetCollection(ctx, in.CollectionID)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("get collection from repository: %w", err)
	}
	if collection.CustomerID != in.CustomerID {
		return entity.Collection{}, apperror.ErrCollectionAccessDenied
	}

	updatedCollection, err := s.collectionRepo.UpdateCollection(ctx, in)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("update collection in repository: %w", err)
	}

	return updatedCollection, nil
}
