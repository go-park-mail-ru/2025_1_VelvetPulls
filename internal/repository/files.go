package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type IFilesRepo interface {
	GetFile(ctx context.Context, fileID string, userID string) (*bytes.Buffer, *model.FileMetaData, error)
	SaveFile(ctx context.Context, fileBuffer *bytes.Buffer, metadata model.FileMetaData, allowedUsers []string) (string, error)
	DeleteFile(ctx context.Context, fileID string, userID string) error
	RewriteFile(ctx context.Context, fileID string, fileBuffer *bytes.Buffer, metadata model.FileMetaData) error

	CreateSticker(ctx context.Context, fileBuffer *bytes.Buffer, metadata model.FileMetaData, packName string) error
	GetStickerPack(ctx context.Context, packID string) (model.StickerPack, error)
	GetStickerPacks(ctx context.Context) ([]model.StickerPack, error)
}

type filesRepository struct {
	minioClient *minio.Client
	bucketName  string
}

func NewFilesRepo(minioClient *minio.Client, bucketName string) IFilesRepo {
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil
		}
	}
	return &filesRepository{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (r *filesRepository) GetFile(ctx context.Context, fileID, userID string) (*bytes.Buffer, *model.FileMetaData, error) {
	info, err := r.minioClient.StatObject(ctx, r.bucketName, fileID, minio.StatObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	if err := checkAccess(info, userID); err != nil {
		return nil, nil, err
	}

	obj, err := r.minioClient.GetObject(ctx, r.bucketName, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, obj); err != nil {
		return nil, nil, err
	}

	size, _ := strconv.ParseInt(info.UserMetadata["X-Amz-Meta-Size"], 10, 64)
	fileMeta := &model.FileMetaData{
		Filename:    info.UserMetadata["X-Amz-Meta-Filename"],
		ContentType: info.UserMetadata["X-Amz-Meta-Content-Type"],
		FileSize:    size,
	}
	return buf, fileMeta, nil
}

func (r *filesRepository) SaveFile(ctx context.Context, buf *bytes.Buffer, meta model.FileMetaData, allowedUsers []string) (string, error) {
	fileID := generateID()
	userMeta := map[string]string{
		"filename":     meta.Filename,
		"content-type": meta.ContentType,
		"size":         strconv.FormatInt(meta.FileSize, 10),
		"users":        strings.Join(allowedUsers, ","),
	}
	_, err := r.minioClient.PutObject(ctx, r.bucketName, fileID,
		bytes.NewReader(buf.Bytes()), int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType:  meta.ContentType,
			UserMetadata: userMeta,
		},
	)
	return fileID, err
}

func (r *filesRepository) DeleteFile(ctx context.Context, fileID string, userID string) error {
	info, err := r.minioClient.StatObject(ctx, r.bucketName, fileID, minio.StatObjectOptions{})
	if err != nil {
		return err
	}
	if err := checkAccess(info, userID); err != nil {
		return err
	}
	return r.minioClient.RemoveObject(ctx, r.bucketName, fileID, minio.RemoveObjectOptions{})
}

func (r *filesRepository) RewriteFile(ctx context.Context, fileID string, fileBuffer *bytes.Buffer, metadata model.FileMetaData) error {
	if fileBuffer == nil || fileBuffer.Len() == 0 {
		return errors.New("file buffer is empty")
	}

	// Получаем текущую информацию об объекте, чтобы сохранить users
	info, err := r.minioClient.StatObject(ctx, r.bucketName, fileID, minio.StatObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to stat object before rewrite: %w", err)
	}

	// Сохраняем старый список пользователей
	users := info.UserMetadata["X-Amz-Meta-Users"]

	// Перезаписываем файл с новыми метаданными
	userMeta := map[string]string{
		"filename":     metadata.Filename,
		"content-type": metadata.ContentType,
		"size":         strconv.FormatInt(metadata.FileSize, 10),
		"users":        users,
	}

	_, err = r.minioClient.PutObject(ctx, r.bucketName, fileID, bytes.NewReader(fileBuffer.Bytes()), int64(fileBuffer.Len()), minio.PutObjectOptions{
		ContentType:  metadata.ContentType,
		UserMetadata: userMeta,
	})
	if err != nil {
		return fmt.Errorf("failed to rewrite file in minio: %w", err)
	}
	return nil
}

// CreateSticker сохраняет стикер в папку с названием стикерпакa
func (r *filesRepository) CreateSticker(ctx context.Context, fileBuffer *bytes.Buffer, metadata model.FileMetaData, packName string) error {
	// Генерируем уникальное имя файла для стикера
	fileID := fmt.Sprintf("stickers/%s/%s-%d", packName, metadata.Filename, time.Now().UnixNano())

	_, err := r.minioClient.PutObject(ctx, r.bucketName, fileID, bytes.NewReader(fileBuffer.Bytes()), int64(fileBuffer.Len()), minio.PutObjectOptions{
		ContentType: metadata.ContentType,
	})
	if err != nil {
		return err
	}

	// Обновляем/создаем метаинформацию по стикерпаку
	return r.updateStickerPackMeta(ctx, packName, fileID)
}

// updateStickerPackMeta добавляет/обновляет JSON с информацией о стикерпаке
func (r *filesRepository) updateStickerPackMeta(ctx context.Context, packName string, newPhoto string) error {
	// Сначала пытаемся получить существующий JSON с метаинфой
	metaObjectName := fmt.Sprintf("stickers/%s/meta.json", packName)

	object, err := r.minioClient.GetObject(ctx, r.bucketName, metaObjectName, minio.GetObjectOptions{})
	if err != nil {
		// Если объекта нет — создаём новый
		return r.saveStickerPackMeta(ctx, packName, newPhoto)
	}
	defer object.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object)
	if err != nil {
		return err
	}

	var pack model.StickerPack
	err = json.Unmarshal(buf.Bytes(), &pack)
	if err != nil {
		// Если не удалось распарсить — создаём новый
		return r.saveStickerPackMeta(ctx, packName, newPhoto)
	}

	// Обновляем photo и packID (если пусто)
	pack.Photo = newPhoto
	if pack.PackID == uuid.Nil {
		pack.PackID = uuid.New()
	}

	return r.saveStickerPackMeta(ctx, packName, pack.Photo)
}

// saveStickerPackMeta сохраняет JSON с метаинформацией
func (r *filesRepository) saveStickerPackMeta(ctx context.Context, packName, photo string) error {
	metaObjectName := fmt.Sprintf("stickers/%s/meta.json", packName)
	pack := model.StickerPack{
		Photo:  photo,
		PackID: uuid.New(),
	}
	metaBytes, err := json.Marshal(pack)
	if err != nil {
		return err
	}

	_, err = r.minioClient.PutObject(ctx, r.bucketName, metaObjectName, bytes.NewReader(metaBytes), int64(len(metaBytes)), minio.PutObjectOptions{
		ContentType: "application/json",
	})
	return err
}

// GetStickerPack возвращает стикерпак по packID (ищет в бакете все meta.json и сравнивает)
func (r *filesRepository) GetStickerPack(ctx context.Context, packID string) (model.StickerPack, error) {
	// Перебираем все объекты с префиксом stickers/
	// Ищем meta.json файлов и сравниваем packID
	var foundPack model.StickerPack

	opts := minio.ListObjectsOptions{
		Prefix:    "stickers/",
		Recursive: true,
	}

	for object := range r.minioClient.ListObjects(ctx, r.bucketName, opts) {
		if object.Err != nil {
			continue
		}
		if !endsWith(object.Key, "meta.json") {
			continue
		}

		obj, err := r.minioClient.GetObject(ctx, r.bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			continue
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, obj)
		obj.Close()
		if err != nil {
			continue
		}

		var pack model.StickerPack
		err = json.Unmarshal(buf.Bytes(), &pack)
		if err != nil {
			continue
		}

		if pack.PackID.String() == packID {
			foundPack = pack
			return foundPack, nil
		}
	}

	return model.StickerPack{}, fmt.Errorf("sticker pack not found")
}

// GetStickerPacks возвращает список всех стикерпаков
func (r *filesRepository) GetStickerPacks(ctx context.Context) ([]model.StickerPack, error) {
	var packs []model.StickerPack

	opts := minio.ListObjectsOptions{
		Prefix:    "stickers/",
		Recursive: true,
	}

	for object := range r.minioClient.ListObjects(ctx, r.bucketName, opts) {
		if object.Err != nil {
			continue
		}
		if !endsWith(object.Key, "meta.json") {
			continue
		}

		obj, err := r.minioClient.GetObject(ctx, r.bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			continue
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, obj)
		obj.Close()
		if err != nil {
			continue
		}

		var pack model.StickerPack
		err = json.Unmarshal(buf.Bytes(), &pack)
		if err != nil {
			continue
		}

		packs = append(packs, pack)
	}

	return packs, nil
}

// простая функция для проверки окончания строки
func endsWith(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}

func checkAccess(info minio.ObjectInfo, userID string) error {
	usersMeta := info.UserMetadata["X-Amz-Meta-Users"]
	if usersMeta == "" {
		// нет метаданных — либо публичный, либо запрещаем всем
		return nil
	}
	for _, uid := range strings.Split(usersMeta, ",") {
		if strings.TrimSpace(uid) == userID {
			return nil
		}
	}
	return errors.New("forbidden")
}

func generateID() string {
	return uuid.New().String()
}
