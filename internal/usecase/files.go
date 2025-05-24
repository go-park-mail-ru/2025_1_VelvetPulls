package usecase

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type filesUsecase struct {
	fileRepo repository.IFilesRepo
}

type IFilesUsecase interface {
	GetFile(ctx context.Context, fileIDStr uuid.UUID, userID uuid.UUID) (*bytes.Buffer, *model.FileMetaData, error)
	SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (model.Payload, error)
	GetStickerPack(ctx context.Context, packID string) (model.GetStickerPackResponse, error)
	GetStickerPacks(ctx context.Context) (model.StickerPacks, error)
	SaveSticker(ctx context.Context, file multipart.File, header *multipart.FileHeader, name string) error
	SavePhoto(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (model.Payload, error)
	// SaveAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	// RewritePhoto(ctx context.Context, file multipart.File, header multipart.FileHeader, fileIDStr string) error
	// DeletePhoto(ctx context.Context, fileIDStr string) error
	// UpdateFile(ctx context.Context, fileIDStr string, file multipart.File, header *multipart.FileHeader) (string, error)
}

func NewFilesUsecase(fileRepo repository.IFilesRepo) IFilesUsecase {
	return &filesUsecase{fileRepo: fileRepo}
}

func (u *filesUsecase) GetFile(ctx context.Context, fileIDStr uuid.UUID, userID uuid.UUID) (*bytes.Buffer, *model.FileMetaData, error) {
	return u.fileRepo.GetFile(ctx, fileIDStr, userID)
}

func (u *filesUsecase) SaveFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (model.Payload, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	logger.Info("Starting file save operation",
		zap.String("filename", header.Filename),
		zap.Int64("size", header.Size),
		zap.Int("users_count", len(users)),
	)

	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		logger.Error("Failed to create file buffer",
			zap.String("filename", header.Filename),
			zap.Error(err),
		)
		fmt.Println(err)
		return model.Payload{}, fmt.Errorf("failed to create file buffer: %w", err)
	}

	fileID, err := u.fileRepo.SaveFile(
		ctx,
		fileBuffer,
		header.Filename,
		header.Header.Get("Content-Type"),
		header.Size,
		users,
	)
	if err != nil {
		fmt.Println(err)
		logger.Error("Failed to save file in repository",
			zap.String("filename", header.Filename),
			zap.Error(err),
		)
		return model.Payload{}, fmt.Errorf("failed to save file in repository: %w", err)
	}

	out := model.Payload{
		URL:         addFileURLPrefix(fileID),
		Filename:    header.Filename,
		ContentType: "file",
		Size:        header.Size,
	}

	logger.Info("File successfully saved",
		zap.String("file_id", fileID),
		zap.String("result_url", out.URL),
	)

	return out, nil
}

func (u *filesUsecase) SavePhoto(ctx context.Context, file multipart.File, header *multipart.FileHeader, users []string) (model.Payload, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	logger.Info("Starting photo save operation",
		zap.String("filename", header.Filename),
		zap.Int64("size", header.Size),
		zap.Int("users_count", len(users)),
	)

	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		logger.Error("Failed to create photo buffer",
			zap.String("filename", header.Filename),
			zap.Error(err),
		)
		return model.Payload{}, fmt.Errorf("failed to create photo buffer: %w", err)
	}

	fileID, err := u.fileRepo.SaveFile(
		ctx,
		fileBuffer,
		header.Filename,
		header.Header.Get("Content-Type"),
		header.Size,
		users,
	)
	if err != nil {
		logger.Error("Failed to save photo in repository",
			zap.String("filename", header.Filename),
			zap.Error(err),
		)
		return model.Payload{}, fmt.Errorf("failed to save photo in repository: %w", err)
	}

	out := model.Payload{
		URL:         addFileURLPrefix(fileID),
		Filename:    header.Filename,
		ContentType: "photo", // указываем явно, что это фото
		Size:        header.Size,
	}

	logger.Info("Photo successfully saved",
		zap.String("file_id", fileID),
		zap.String("result_url", out.URL),
	)

	return out, nil
}

func (u *filesUsecase) GetStickerPack(ctx context.Context, packID string) (model.GetStickerPackResponse, error) {
	return u.fileRepo.GetStickerPack(ctx, packID)
}

func (u *filesUsecase) GetStickerPacks(ctx context.Context) (model.StickerPacks, error) {
	return u.fileRepo.GetStickerPacks(ctx)
}

func (u *filesUsecase) SaveSticker(ctx context.Context, file multipart.File, header *multipart.FileHeader, packName string) error {
	if file == nil || header == nil {
		return fmt.Errorf("invalid file data")
	}
	fileBuffer, err := getFileBuffer(file)
	if err != nil {
		return fmt.Errorf("failed to create file buffer: %w", err)
	}

	metadata := model.FileMetaData{
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		FileSize:    header.Size,
	}
	if _, err := u.fileRepo.CreateSticker(ctx, fileBuffer, metadata, packName); err != nil {
		return fmt.Errorf("failed to save sticker: %w", err)
	}

	return nil
}

func getFileBuffer(file multipart.File) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, err
	}
	return buf, nil
}

func addFileURLPrefix(fileID string) string {
	return "/files/" + fileID
}
