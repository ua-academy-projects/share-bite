package collection

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
)

const (
	maxLimit     = 100
	defaultLimit = 20
)

func (s *service) ListCustomerCollections(
	ctx context.Context,
	in entity.ListCustomerCollectionsInput,
) (entity.ListCustomerCollectionsOutput, error) {
	cursorTime, cursorID, err := s.parsePageToken(in.PageToken)
	if err != nil {
		return entity.ListCustomerCollectionsOutput{}, apperror.ErrInvalidPageToken
	}

	limit := in.PageSize
	switch {
	case limit <= 0:
		limit = defaultLimit
	case limit > maxLimit:
		limit = maxLimit
	}

	// we always ask to return limit + 1
	collections, err := s.collectionRepo.ListCustomerCollections(ctx, in.CustomerID, cursorTime, cursorID, limit+1)
	if err != nil {
		return entity.ListCustomerCollectionsOutput{}, fmt.Errorf("get customer's collections from repository: %w", err)
	}

	var nextPageToken string

	// if len(rows) > limit -> 1 more page exists
	if len(collections) > limit {
		collections = collections[:limit]

		lastItem := collections[len(collections)-1]
		nextPageToken = s.generatePageToken(lastItem.CreatedAt, lastItem.ID)
	}

	return entity.ListCustomerCollectionsOutput{
		Collections:   collections,
		NextPageToken: nextPageToken,
	}, nil
}
