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
	if err := s.requireCollaborator(ctx, in.CollectionID, in.CustomerID, collection.CustomerID); err != nil {
		return entity.Collection{}, err
	}

	// only owner can have control over collection visibility
	if in.IsPublic != nil && in.CustomerID != collection.CustomerID {
		return entity.Collection{}, apperror.ErrCollectionAccessDenied
	}

	updatedCollection, err := s.collectionRepo.UpdateCollection(ctx, in)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("update collection in repository: %w", err)
	}

	return updatedCollection, nil
}
