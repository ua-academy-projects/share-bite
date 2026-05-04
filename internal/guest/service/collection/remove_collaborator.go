package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) RemoveCollaborator(ctx context.Context, in entity.RemoveCollaboratorInput) error {
	collection, err := s.collectionRepo.GetCollection(ctx, in.CollectionID)
	if err != nil {
		return fmt.Errorf("get collection from repository: %w", err)
	}
	if err := s.requireOwner(ctx, in.CollectionID, in.CustomerID, collection.CustomerID); err != nil {
		return err
	}

	if err := s.collectionRepo.DeleteCollaborator(ctx, in.CollectionID, in.TargetCustomerID); err != nil {
		return fmt.Errorf("delete collaborator in repository: %w", err)
	}

	return nil
}
