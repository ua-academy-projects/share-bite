package storage

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Upload(context.Context, string, string, io.Reader) (string, error)
	Delete(context.Context, string) error
	BuildURL(string) string
}
