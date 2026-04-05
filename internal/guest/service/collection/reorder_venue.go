package collection

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const (
	rebalanceGapLimit = 1e-9
)

func (s *service) ReorderVenue(ctx context.Context, in entity.ReorderVenueInput) error {
	collection, err := s.collectionRepo.GetCollection(ctx, in.CollectionID)
	if err != nil {
		return fmt.Errorf("get collection from repository: %w", err)
	}
	if collection.CustomerID != in.CustomerID {
		return apperror.ErrCollectionAccessDenied
	}

	exists, err := s.collectionRepo.CheckIfVenueInCollection(ctx, in.CollectionID, in.VenueID)
	if err != nil {
		return fmt.Errorf("check if venue is in collection from repository: %w", err)
	}
	if !exists {
		return apperror.VenueNotFoundInCollection(in.VenueID)
	}

	newSortOrder, gap, err := s.generateNewSortOrder(ctx, in)
	if err != nil {
		return fmt.Errorf("generate new sort order: %w", err)
	}

	if err := s.collectionRepo.UpdateVenueSortOrder(
		ctx,
		in.CollectionID,
		in.VenueID,
		newSortOrder,
	); err != nil {
		return fmt.Errorf("update venue sort order in repository: %w", err)
	}

	if gap > 0 && gap < rebalanceGapLimit {
		go func(collectionID string) {
			rebalanceCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			logger.Info(rebalanceCtx, "starting async rebalance for collection: ", collectionID)

			if err := s.collectionRepo.RebalanceCollectionSortOrders(rebalanceCtx, collectionID); err != nil {
				logger.Errorf(rebalanceCtx, "rebalance for collection %q failed: %v", collectionID, err)
				return
			}

			logger.Infof(rebalanceCtx, "rebalance for collection %q successfully completed", collectionID)
		}(in.CollectionID)
	}

	return nil
}

func (s *service) generateNewSortOrder(ctx context.Context, in entity.ReorderVenueInput) (float64, float64, error) {
	if in.PrevVenueID == nil && in.NextVenueID == nil {
		return 0, 0, apperror.ErrInvalidReorderParams
	}

	var (
		orderAbove float64
		orderBelow float64
	)

	if in.PrevVenueID != nil {
		prevVenue, err := s.collectionRepo.GetCollectionVenue(ctx, in.CollectionID, *in.PrevVenueID)
		if err != nil {
			return 0, 0, fmt.Errorf("get prev venue from collection from repository: %w", err)
		}

		orderAbove = prevVenue.SortOrder
	}
	if in.NextVenueID != nil {
		nextVenue, err := s.collectionRepo.GetCollectionVenue(ctx, in.CollectionID, *in.NextVenueID)
		if err != nil {
			return 0, 0, fmt.Errorf("get next venue from collection from repository: %w", err)
		}

		orderBelow = nextVenue.SortOrder
	}

	var (
		newSortOrder float64
		currGap      float64
	)

	if in.PrevVenueID != nil && in.NextVenueID != nil {
		// between two venues
		newSortOrder = (orderAbove + orderBelow) / 2.0
		currGap = orderBelow - orderAbove

	} else if in.PrevVenueID == nil && in.NextVenueID != nil {
		// no venue above
		newSortOrder = orderBelow / 2.0
		currGap = orderBelow

	} else if in.PrevVenueID != nil && in.NextVenueID == nil {
		// the least in the list
		newSortOrder = orderAbove + sortOrderGap
		currGap = sortOrderGap
	}

	return newSortOrder, currGap, nil
}
