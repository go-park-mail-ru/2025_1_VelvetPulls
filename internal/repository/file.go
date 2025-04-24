package repository

import (
	"context"
	"io"
)

type FileRepository interface {
	GetFile(ctx context.Context, objectPath string) (io.ReadCloser, error)
	SaveFile(ctx context.Context, objectPath string, file io.Reader, size int64, contentType string) error
	DeleteFile(ctx context.Context, objectPath string) error
}
