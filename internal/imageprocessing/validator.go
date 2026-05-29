package imageprocessing

import "fmt"

const (
	MinImageWidth  = 200
	MinImageHeight = 200

	MaxImageWidth  = 10000
	MaxImageHeight = 10000
)

func ValidateDimensions(width, height int) error {
	if width < MinImageWidth || height < MinImageHeight {
		return fmt.Errorf(
			"image dimensions too small: %dx%d",
			width,
			height,
		)
	}

	if width > MaxImageWidth || height > MaxImageHeight {
		return fmt.Errorf(
			"image dimensions too large: %dx%d",
			width,
			height,
		)
	}

	return nil
}
