package business

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/internal/storage/mediatype"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

const (
	kilometerIndex = 1.60934
)

func (s *service) CreateBox(ctx context.Context, userID string, req dto.CreateBoxRequest, image *multipart.FileHeader) (*entity.Box, error) {
	const op = "service.box.CreateBox"

	if image == nil {
		return nil, fmt.Errorf("%s: image is required", op)
	}

	openedFile, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	buffer := make([]byte, fileHeaderSize)
	n, err := openedFile.Read(buffer)
	openedFile.Close()

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	contentType := http.DetectContentType(buffer[:n])
	if err := postImageValidator.Validate(contentType, image.Size); err != nil {
		if errors.Is(err, mediatype.ErrUnsupportedType) {
			return nil, biserr.WrongFileExtErr
		}

		if errors.Is(err, mediatype.ErrFileTooLarge) {
			return nil, biserr.FileToLargeErr
		}

		return nil, fmt.Errorf("%s: validation failed: %w", op, err)
	}

	if req.DiscountPrice.GreaterThan(req.FullPrice) {
		return nil, fmt.Errorf("%s: %w", op, errors.New("invalid price"))
	}

	if req.FullPrice.LessThanOrEqual(decimal.Zero) || req.DiscountPrice.LessThan(decimal.Zero) {
		return nil, fmt.Errorf("%s: %w", op, errors.New("price values are out of range"))
	}

	if req.Quantity <= 0 {
		return nil, fmt.Errorf("%s: %w", op, errors.New("quantity must be at least 1"))
	}

	openedFile, err = image.Open()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fileExt := filepath.Ext(image.Filename)
	objectKey := fmt.Sprintf("boxes/%d/%s%s", req.VenueID, uuid.New().String(), fileExt)

	err = s.storage.Upload(ctx, objectKey, contentType, openedFile)
	openedFile.Close()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	isSuccess := false
	defer func() {
		if !isSuccess {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := s.storage.Delete(cleanupCtx, objectKey); err != nil {
				log.Printf("failed to cleanup uploaded object: key=%s err=%v", objectKey, err)
			}
		}
	}()

	var box *entity.Box

	err = s.txManager.ReadCommitted(ctx, func(ctxTx context.Context) error {
		err := s.businessRepo.CheckOwnership(ctxTx, userID, req.VenueID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		box = &entity.Box{
			VenueID:       req.VenueID,
			CategoryID:    req.CategoryID,
			Image:         objectKey,
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

	isSuccess = true
	return box, nil
}

// func (s *service) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int, orgID *int) (pagination.Result[entity.BoxWithDistance], error) {
// 	const op = "service.box.ListNearbyBoxes"

// 	result, err := s.businessRepo.ListNearbyBoxes(ctx, offset, limit, lat, lon, categoryID, orgID)
// 	if err != nil {
// 		return pagination.Result[entity.BoxWithDistance]{}, fmt.Errorf("%s: %w", op, err)
// 	}

// 	for i := range result.Items {
// 		result.Items[i].Distance = result.Items[i].Distance * kilometerIndex

// 		if result.Items[i].Box.Image != "" {
// 			imageURL, err := s.storage.GetPresignedURL(ctx, result.Items[i].Box.Image)
// 			if err != nil {
// 				log.Printf(
// 					"failed to generate presigned URL for box %d: %v",
// 					result.Items[i].Box.ID,
// 					err,
// 				)

// 				result.Items[i].Box.Image = ""
// 			} else {
// 				result.Items[i].Box.Image = imageURL
// 			}
// 		}
// 	}

// 	return result, nil
// }

func (s *service) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int, orgID *int) (pagination.Result[entity.BoxWithDistance], error) {
    const op = "service.box.ListNearbyBoxes"

    result, err := s.businessRepo.ListNearbyBoxes(ctx, offset, limit, lat, lon, categoryID, orgID)
    if err != nil {
        return pagination.Result[entity.BoxWithDistance]{}, fmt.Errorf("%s: %w", op, err)
    }

    for i := range result.Items {
        result.Items[i].Distance = result.Items[i].Distance * kilometerIndex

        img := result.Items[i].Box.Image
        if img != "" {
            if strings.HasPrefix(img, "http://") || strings.HasPrefix(img, "https://") {
                result.Items[i].Box.Image = img
            } else {
                result.Items[i].Box.Image = s.storage.BuildURL(img)
            }
        }
    }

    return result, nil
}

func (s *service) ReserveBox(ctx context.Context, userID string, boxID int64) (*entity.BoxReservation, error) {
	const op = "service.box.ReserveBox"

	box, err := s.businessRepo.GetBox(ctx, boxID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	boxCode, err := s.businessRepo.ReserveBoxItem(ctx, boxID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entity.BoxReservation{
		Image:         box.Image,
		FullPrice:     box.FullPrice,
		DiscountPrice: box.DiscountPrice,
		BoxCode:       boxCode,
	}, nil
}
