package imageprocessing

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/uuid"

	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/internal/storage/key"
)

type Repository interface {
	UpdateProcessedMetadata(ctx context.Context, imageID string, thumbnailKey string, width int, height int) error
	MarkProcessingFailed(ctx context.Context, imageID string, reason string) error
	IsAlreadyProcessed(ctx context.Context, imageID string) (bool, error)
	MarkProcessing(ctx context.Context, imageID string) error
}

type Service struct {
	storage    storage.ObjectStorage
	repository Repository
}

func NewService(storage storage.ObjectStorage, repository Repository) *Service {
	return &Service{
		storage:    storage,
		repository: repository,
	}
}

func (s *Service) ProcessImage(ctx context.Context, msg ProcessImageMessage) (err error) {
	defer func() {
		if err != nil {
			_ = s.repository.MarkProcessingFailed(ctx, msg.ImageID, err.Error())
		}
	}()

	alreadyProcessed, err := s.repository.IsAlreadyProcessed(ctx, msg.ImageID)
	if err != nil {
		return err
	}

	if alreadyProcessed {
		return nil
	}

	err = s.repository.MarkProcessing(ctx, msg.ImageID)
	if err != nil {
		return err
	}

	reader, err := s.storage.Get(ctx, msg.S3Key)
	if err != nil {
		return err
	}

	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	bounds := img.Bounds()

	width := bounds.Dx()
	height := bounds.Dy()

	if err := ValidateDimensions(width, height); err != nil {
		return err
	}

	thumbnailBuffer, err := GenerateThumbnail(img)
	if err != nil {
		return err
	}

	thumbnailKey := key.PostThumbnailKey(uuid.NewString(), uuid.NewString(), "jpg")

	err = s.storage.Upload(ctx, thumbnailKey, "image/jpeg", bytes.NewReader(thumbnailBuffer.Bytes()))
	if err != nil {
		return err
	}

	err = s.repository.UpdateProcessedMetadata(ctx, msg.ImageID, thumbnailKey, width, height)
	if err != nil {
		return err
	}

	return nil
}
