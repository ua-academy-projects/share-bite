package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) UpdateCollection(ctx context.Context, in entity.UpdateCollectionInput) (entity.Collection, error) {
	var updatedCollection entity.Collection

	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		collection, err := s.collectionRepo.GetCollectionForUpdate(ctx, in.CollectionID)
		if err != nil {
			return fmt.Errorf("get collection from repository: %w", err)
		}
		if err := s.requireCollaborator(ctx, in.CollectionID, in.CustomerID, collection.CustomerID); err != nil {
			return err
		}

		// only owner can have control over collection visibility
		if in.IsPublic != nil && in.CustomerID != collection.CustomerID {
			return apperror.ErrCollectionAccessDenied
		}

		updatedCollection, err = s.collectionRepo.UpdateCollection(ctx, in)
		if err != nil {
			return fmt.Errorf("update collection in repository: %w", err)
		}

		return nil
	}); txErr != nil {
		return entity.Collection{}, txErr
	}

	return updatedCollection, nil
}
