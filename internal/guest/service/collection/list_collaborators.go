package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

func (s *service) ListCollaborators(ctx context.Context, collectionID string, customerID *string) ([]entity.Collaborator, error) {
	collection, err := s.collectionRepo.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get collection from repository: %w", err)
	}

	if !collection.IsPublic {
		if customerID == nil {
			return nil, apperror.CollectionNotFoundID(collectionID)
		}

		if err := s.requireCollaborator(ctx, collectionID, *customerID, collection.CustomerID); err != nil {
			return nil, err
		}
	}

	collaborators, err := s.collectionRepo.ListCollaborators(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("get collaborators from repository: %w", err)
	}

	return collaborators, nil
}
