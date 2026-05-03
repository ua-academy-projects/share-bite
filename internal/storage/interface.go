package storage

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Upload(context.Context, string, string, io.Reader) error
	Delete(context.Context, string) error
	BuildURL(string) string //TODO: delete with complete AWS integration
	GetPresignedURL(ctx context.Context, key string) (string, error)
}
