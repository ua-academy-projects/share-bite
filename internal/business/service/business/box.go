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

func (s *service) ListBoxesByBusiness(ctx context.Context, userID string, offset, limit int) (pagination.Result[entity.Box], error) {
	const op = "service.box.ListBoxesByBusiness"

	brandID, err := s.businessRepo.GetBrandIDByOwnerUserID(ctx, userID)
	if err != nil {
		return pagination.Result[entity.Box]{}, fmt.Errorf("%s: %w", op, err)
	}

	orgs, err := s.businessRepo.ListByParentID(ctx, brandID, 0, 1000, nil)
	if err != nil {
		return pagination.Result[entity.Box]{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(orgs.Items) == 0 {
		return pagination.Result[entity.Box]{Items: []entity.Box{}, Total: 0}, nil
	}

	venueID := orgs.Items[0].Id
	result, err := s.businessRepo.ListBoxesByVenueID(ctx, venueID, offset, limit)
	if err != nil {
		return pagination.Result[entity.Box]{}, fmt.Errorf("%s: %w", op, err)
	}

	for i := range result.Items {
		if result.Items[i].Image != "" {
			if !strings.HasPrefix(result.Items[i].Image, "http://") && !strings.HasPrefix(result.Items[i].Image, "https://") {
				result.Items[i].Image = s.storage.BuildURL(result.Items[i].Image)
			}
		}
	}

	return result, nil
}

func (s *service) UpdateBox(ctx context.Context, boxID int64, userID string, input entity.BoxUpdateInput) (*entity.Box, error) {
	const op = "service.box.UpdateBox"

	box, err := s.businessRepo.GetBox(ctx, boxID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.businessRepo.CheckOwnership(ctx, userID, box.VenueID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if input.FullPrice != nil && input.FullPrice.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("%s: full price must be greater than 0", op)
	}

	if input.DiscountPrice != nil && input.DiscountPrice.LessThan(decimal.Zero) {
		return nil, fmt.Errorf("%s: discount price cannot be negative", op)
	}

	fullPrice := box.FullPrice
	discountPrice := box.DiscountPrice

	if input.FullPrice != nil {
		fullPrice = *input.FullPrice
	}

	if input.DiscountPrice != nil {
		discountPrice = *input.DiscountPrice
	}

	if discountPrice.GreaterThan(fullPrice) {
		return nil, fmt.Errorf("%s: discount price cannot be greater than full price", op)
	}

	if input.ExpiresAt != nil && !input.ExpiresAt.After(time.Now()) {
		return nil, fmt.Errorf("%s: expires_at must be in the future", op)
	}

	updatedBox, err := s.businessRepo.UpdateBox(ctx, boxID, input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if updatedBox.Image != "" {
		if !strings.HasPrefix(updatedBox.Image, "http://") && !strings.HasPrefix(updatedBox.Image, "https://") {
			updatedBox.Image = s.storage.BuildURL(updatedBox.Image)
		}
	}

	return updatedBox, nil
}

func (s *service) GetBoxReservations(ctx context.Context, boxID int64, userID string, offset, limit int) (pagination.Result[entity.BoxItem], error) {
	const op = "service.box.GetBoxReservations"

	box, err := s.businessRepo.GetBox(ctx, boxID)
	if err != nil {
		return pagination.Result[entity.BoxItem]{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.businessRepo.CheckOwnership(ctx, userID, box.VenueID)
	if err != nil {
		return pagination.Result[entity.BoxItem]{}, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.businessRepo.GetBoxItems(ctx, boxID, offset, limit)
	if err != nil {
		return pagination.Result[entity.BoxItem]{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
