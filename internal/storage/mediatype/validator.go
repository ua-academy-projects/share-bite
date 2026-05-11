package mediatype

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedType = errors.New("unsupported image type")
	ErrFileTooLarge    = errors.New("file size exceeds the limit")
)

const (
	// DefaultMaxImageSizeBytes is the default maximum allowed image upload size.
	DefaultMaxImageSizeBytes = 5 * 1024 * 1024 // 5MB
)

// ImageValidator validates uploaded image files by MIME type and size.
type ImageValidator struct {
	allowedMIMETypes map[string]struct{}
	maxSizeBytes     int64
}

var (
	// DefaultImageValidator is a pre-configured validator for standard image uploads.
	DefaultImageValidator = &ImageValidator{
		allowedMIMETypes: map[string]struct{}{
			"image/jpeg": {},
			"image/png":  {},
		},
		maxSizeBytes: DefaultMaxImageSizeBytes,
	}
)

// NewValidator creates a custom ImageValidator with the given size limit and allowed MIME types.
func NewValidator(maxSizeBytes int64, allowedMIMETypes ...string) *ImageValidator {
	v := &ImageValidator{
		maxSizeBytes: maxSizeBytes,
	}

	if len(allowedMIMETypes) > 0 {
		mimeTypes := make(map[string]struct{}, len(allowedMIMETypes))
		for _, m := range allowedMIMETypes {
			mimeTypes[m] = struct{}{}
		}

		v.allowedMIMETypes = mimeTypes
	}

	return v
}

// Validate checks whether the given content type is allowed and the file size is within the limit.
func (v *ImageValidator) Validate(contentType string, sizeBytes int64) error {
	if _, ok := v.allowedMIMETypes[contentType]; !ok {
		return ErrUnsupportedType
	}

	if sizeBytes > v.maxSizeBytes {
		return fmt.Errorf("%w: max %dMB, got %d bytes",
			ErrFileTooLarge,
			v.maxSizeBytes/1024/1024,
			sizeBytes,
		)
	}

	return nil
}
