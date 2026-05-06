package collection

import (
	"context"
	"fmt"
)

func (s *service) DeleteCollection(ctx context.Context, collectionID string, customerID string) error {
	collection, err := s.collectionRepo.GetCollection(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("get collection from repository: %w", err)
	}
	if err := s.requireOwner(ctx, collectionID, customerID, collection.CustomerID); err != nil {
		return err
	}

	if err := s.collectionRepo.DeleteCollection(ctx, collectionID); err != nil {
		return fmt.Errorf("delete collection in repository: %w", err)
	}

	return nil
}
