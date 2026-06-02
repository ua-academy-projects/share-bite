package imageprocessing

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

const (
	ThumbnailWidth  = 300
	ThumbnailHeight = 300
)

func GenerateThumbnail(img image.Image) (*bytes.Buffer, error) {
	thumbnail := resize.Thumbnail(
		ThumbnailWidth,
		ThumbnailHeight,
		img,
		resize.Lanczos3,
	)

	var buffer bytes.Buffer

	err := jpeg.Encode(
		&buffer,
		thumbnail,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
