package addPhoto

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"io"
	"path/filepath"
	"strings"
	"time"
)

type UploadPhotoInput struct {
	File        io.Reader
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type UploadPhotoOutput struct {
	ObjectKey string `json:"object_key"`
	ImageURL  string `json:"image_url"`
}

type Service struct {
	storage storage.ObjectStorage
}

func NewService(storage storage.ObjectStorage) *Service {
	return &Service{storage: storage}
}

func (s *Service) UploadPhoto(ctx context.Context, in UploadPhotoInput) (*UploadPhotoOutput, error) {
	if in.File == nil {
		return nil, apperror.BadRequest("missing file to upload")
	}
	if in.Filename == "" {
		return nil, apperror.BadRequest("filename is required")
	}
	if in.ContentType == "" {
		return nil, apperror.BadRequest("content type is required")
	}
	if in.Size <= 0 {
		return nil, apperror.BadRequest("missing size to upload")
	}
	if !isAllowedContentType(in.ContentType) {
		return nil, apperror.BadRequest("content type is not allowed")
	}
	objectKey := generateObjectKey(in.Filename)

	uploadedKey, err := s.storage.Upload(ctx, objectKey, in.ContentType, in.File)
	if err != nil {
		return nil, err
	}

	return &UploadPhotoOutput{
		ObjectKey: uploadedKey,
		ImageURL:  s.storage.BuildURL(uploadedKey),
	}, nil
}

func generateObjectKey(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".bin"
	}
	return fmt.Sprintf("tmp/%d-%s%s", time.Now().Unix(), uuid.NewString(), ext)
}

func isAllowedContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/jpg", "image/png":
		return true
	default:
		return false
	}
}
