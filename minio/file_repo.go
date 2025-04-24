package minio

import (
	"context"
	"io"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/minio/minio-go/v7"
)

type FileRepository struct {
	client *minio.Client
	bucket string
}

func NewFileRepository(client *minio.Client, bucket string) repository.FileRepository {
	return &FileRepository{
		client: client,
		bucket: bucket,
	}
}

func (r *FileRepository) GetFile(ctx context.Context, objectPath string) (io.ReadCloser, error) {
	obj, err := r.client.GetObject(ctx, r.bucket, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, repository.ErrFileStorageFailed
	}

	// Проверяем, существует ли объект
	if _, err = obj.Stat(); err != nil {
		obj.Close()
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, repository.ErrFileNotFound
		}
		return nil, repository.ErrFileStorageFailed
	}

	return obj, nil
}

func (r *FileRepository) SaveFile(ctx context.Context, objectPath string, file io.Reader, size int64, contentType string) error {
	_, err := r.client.PutObject(
		ctx,
		r.bucket,
		objectPath,
		file,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return repository.ErrFileStorageFailed
	}
	return nil
}

func (r *FileRepository) DeleteFile(ctx context.Context, objectPath string) error {
	err := r.client.RemoveObject(
		ctx,
		r.bucket,
		objectPath,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return repository.ErrFileNotFound
		}
		return repository.ErrFileStorageFailed
	}
	return nil
}
