package collection

import (
	"context"
	"fmt"

	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) DeleteCollection(ctx context.Context, collectionID string, customerID string) error {
	collection, err := s.collectionRepo.GetCollection(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("get collection from repository: %w", err)
	}
	if collection.CustomerID != customerID {
		return apperror.ErrCollectionAccessDenied
	}

	if err := s.collectionRepo.DeleteCollection(ctx, collectionID); err != nil {
		return fmt.Errorf("delete collection in repository: %w", err)
	}

	return nil
}
