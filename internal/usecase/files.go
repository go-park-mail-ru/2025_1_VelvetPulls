package usecase

import (
	"bytes"
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
)

type filesUsecase struct {
	fileRepo repository.IFilesRepo
}

type IFilesUsecase interface {
	GetFile(ctx context.Context, fileIDStr string, userID string) (*bytes.Buffer, *model.FileMetaData, error)
	// SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (model.Payload, error)
	// SavePhoto(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (model.Payload, error)
	// SaveAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	// RewritePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader, fileIDStr string) error
	// DeletePhoto(ctx context.Context, fileIDStr string) error
	// UpdateFile(ctx context.Context, fileIDStr string, file multipart.File, header *multipart.FileHeader) (string, error)
	// GetStickerPack(ctx context.Context, packID string) (model.GetStickerPackResponse, error)
}

func NewFilesUsecase(fileRepo repository.IFilesRepo) IFilesUsecase {
	return &filesUsecase{fileRepo: fileRepo}
}

func (u *filesUsecase) GetFile(ctx context.Context, fileIDStr string, userID string) (*bytes.Buffer, *model.FileMetaData, error) {
	return u.fileRepo.GetFile(ctx, fileIDStr, userID)
}
