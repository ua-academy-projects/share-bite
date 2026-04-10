package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func (s *service) CreateBox(ctx context.Context, userID string, req dto.CreateBoxRequest) (*entity.Box, error) {
	const op = "service.box.CreateBox"

	if req.PriceDiscount > req.PriceFull {
		return nil, fmt.Errorf("%s: %w", op, errors.New("invalid price"))
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
			PriceFull:     req.PriceFull,
			PriceDiscount: req.PriceDiscount,
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
