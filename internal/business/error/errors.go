package error

import "errors"

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")

	FileToLargeErr  = errors.New("file is too large")
	WrongFileExtErr = errors.New("unsupported file type")
)