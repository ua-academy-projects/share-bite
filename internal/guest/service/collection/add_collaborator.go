package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

const (
	maxCollaboratorsPerCollection = 100
)

func (s *service) AddCollaborator(ctx context.Context, in entity.AddCollaboratorInput) error {
	if txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		collection, err := s.collectionRepo.GetCollectionForUpdate(ctx, in.CollectionID)
		if err != nil {
			return fmt.Errorf("get collection from repository: %w", err)
		}
		if err := s.requireOwner(ctx, in.CollectionID, in.CustomerID, collection.CustomerID); err != nil {
			return err
		}
		if in.TargetCustomerID == collection.CustomerID {
			return apperror.CustomerAlreadyCollaborator(in.TargetCustomerID)
		}

		isCollaborator, err := s.collectionRepo.CheckIfCollaborator(ctx, in.CollectionID, in.TargetCustomerID)
		if err != nil {
			return fmt.Errorf("check if customer is already a collaborator: %w", err)
		}
		if isCollaborator {
			return apperror.CustomerAlreadyCollaborator(in.TargetCustomerID)
		}

		count, err := s.collectionRepo.CountCollaborators(ctx, in.CollectionID)
		if err != nil {
			return fmt.Errorf("get count of collection collaborators from repository: %w", err)
		}
		if count >= maxCollaboratorsPerCollection {
			return apperror.CollectionCollaboratorsLimitReached(maxCollaboratorsPerCollection)
		}

		if err := s.collectionRepo.CreateCollaborator(ctx, in.CollectionID, in.TargetCustomerID); err != nil {
			return fmt.Errorf("create collaborator in repository: %w", err)
		}

		return nil
	}); txErr != nil {
		return txErr
	}

	return nil
}
