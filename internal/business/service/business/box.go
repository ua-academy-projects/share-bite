package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

const (
	kilometerIndex = 1.60934
)

func (s *service) CreateBox(ctx context.Context, userID string, req dto.CreateBoxRequest) (*entity.Box, error) {
	const op = "service.box.CreateBox"

	if req.DiscountPrice.GreaterThan(req.FullPrice) {
		return nil, fmt.Errorf("%s: %w", op, errors.New("invalid price"))
	}
	if req.FullPrice.LessThanOrEqual(decimal.Zero) || req.DiscountPrice.LessThan(decimal.Zero) {
		return nil, fmt.Errorf("%s: %w", op, errors.New("price values are out of range"))
	}
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("%s: %w", op, errors.New("quantity must be at least 1"))
	}

	var box *entity.Box

	err := s.txManager.ReadCommited(ctx, func(ctxTx context.Context) error {

		err := s.businessRepo.CheckOwnership(ctxTx, userID, req.VenueID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		box = &entity.Box{
			VenueID:       req.VenueID,
			CategoryID:    req.CategoryID,
			Image:         req.Image,
			FullPrice:     req.FullPrice,
			DiscountPrice: req.DiscountPrice,
			ExpiresAt:     req.ExpiresAt,
		}

		boxID, createdAt, err := s.businessRepo.CreateBox(ctxTx, box)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		box.ID = boxID
		box.CreatedAt = createdAt

		for i := 0; i < req.Quantity; i++ {
			code := generateCode()

			err := s.businessRepo.CreateBoxItem(ctxTx, boxID, code)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return box, nil
}

func (s *service) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int) (pagination.Result[entity.BoxWithDistance], error) {
	const op = "service.box.ListNearbyBoxes"

	result, err := s.businessRepo.ListNearbyBoxes(ctx, offset, limit, lat, lon, categoryID)
	if err != nil {
		return pagination.Result[entity.BoxWithDistance]{}, fmt.Errorf("%s: %w", op, err)
	}

	for i := range result.Items {
		result.Items[i].Distance = result.Items[i].Distance * kilometerIndex
	}

	return result, nil
}