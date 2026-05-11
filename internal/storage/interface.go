package storage

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Upload(context.Context, string, string, io.Reader) error
	Delete(context.Context, string) error
	BuildURL(string) string
	GetPresignedURL(ctx context.Context, key string) (string, error)
}
