package collection

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func (s *service) ListCustomerCollections(
	ctx context.Context,
	in entity.ListCustomerCollectionsInput,
) (entity.ListCustomerCollectionsOutput, error) {
	collections, err := s.collectionRepo.ListCustomerCollections(ctx, in.CustomerID, in.CursorTime, in.CursorID, in.Limit)
	if err != nil {
		return entity.ListCustomerCollectionsOutput{}, fmt.Errorf("get customer's collections from repository: %w", err)
	}

	var nextTime *time.Time
	var nextID *string

	requestLimit := in.Limit - 1
	if len(collections) > requestLimit {
		collections = collections[:requestLimit]

		lastItem := collections[len(collections)-1]
		nextTime = &lastItem.CreatedAt
		nextID = &lastItem.ID
	}

	return entity.ListCustomerCollectionsOutput{
		Collections:    collections,
		NextCursorTime: nextTime,
		NextCursorID:   nextID,
	}, nil
}
