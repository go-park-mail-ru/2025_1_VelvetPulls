package repository

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type IFilesRepo interface {
	GetFile(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (*bytes.Buffer, *model.FileMetaData, error)
	SaveFile(ctx context.Context, buf *bytes.Buffer, filename, contentType string, size int64, allowedUsers []string) (string, error)
	DeleteFile(ctx context.Context, fileID string, userID string) error
	RewriteFile(ctx context.Context, fileID string, fileBuffer *bytes.Buffer, metadata model.FileMetaData) error
	CreateSticker(ctx context.Context, fileBuffer *bytes.Buffer, metadata model.FileMetaData, packName string) (uuid.UUID, error)
	GetStickerPack(ctx context.Context, packID string) (model.GetStickerPackResponse, error)
	GetStickerPacks(ctx context.Context) (model.StickerPacks, error)
}

type filesRepository struct {
	minioClient *minio.Client
	db          *sql.DB
	bucketName  string
}

func NewFilesRepo(minioClient *minio.Client, db *sql.DB, bucketName string) IFilesRepo {
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
		db:          db,
		bucketName:  bucketName,
	}
}

func (r *filesRepository) GetFile(ctx context.Context, fileID uuid.UUID, userID uuid.UUID) (*bytes.Buffer, *model.FileMetaData, error) {
	info, err := r.minioClient.StatObject(ctx, r.bucketName, fileID.String(), minio.StatObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	if err := checkAccess(info, userID.String()); err != nil {
		return nil, nil, err
	}
	log.Println(err)

	obj, err := r.minioClient.GetObject(ctx, r.bucketName, fileID.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}
	log.Println(r.bucketName, fileID.String())
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

func (r *filesRepository) SaveFile(ctx context.Context, buf *bytes.Buffer, filename, contentType string, size int64, allowedUsers []string) (string, error) {
	fileID := generateID()

	userMetadata := map[string]string{
		"filename":     filename,
		"content-type": contentType,
		"size":         strconv.FormatInt(size, 10),
	}

	if len(allowedUsers) > 0 {
		userMetadata["allowed-users"] = strings.Join(allowedUsers, ",")
	}

	_, err := r.minioClient.PutObject(
		ctx,
		r.bucketName,
		fileID,
		buf,
		size,
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: userMetadata,
		},
	)
	if err != nil {
		return "", err
	}

	return fileID, nil
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

func (r *filesRepository) CreateSticker(ctx context.Context, fileBuffer *bytes.Buffer, metadata model.FileMetaData, packName string) (uuid.UUID, error) {
	stickerID := uuid.New()

	// Загрузка файла в MinIO (по ID, как имя объекта)
	_, err := r.minioClient.PutObject(ctx, r.bucketName, stickerID.String(), fileBuffer, metadata.FileSize, minio.PutObjectOptions{
		ContentType: metadata.ContentType,
		UserMetadata: map[string]string{
			"X-Amz-Meta-Filename":     metadata.Filename,
			"X-Amz-Meta-Content-Type": metadata.ContentType,
			"X-Amz-Meta-Size":         fmt.Sprintf("%d", metadata.FileSize),
		},
	})
	if err != nil {
		return uuid.Nil, err
	}

	// Добавление в таблицу sticker
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO sticker (id, sticker_path)
		VALUES ($1, $2)
	`, stickerID, fmt.Sprintf("/stickers/%s", stickerID)) // sticker_path — не используется в MinIO, только в UI
	if err != nil {
		return uuid.Nil, err
	}

	// Получение packID, создание при необходимости
	var packID uuid.UUID
	err = r.db.QueryRowContext(ctx, `
		SELECT id FROM sticker_pack WHERE name = $1
	`, packName).Scan(&packID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Пака нет — создаём новую, используем этот стикер как превью
			packID = uuid.New()
			_, err = r.db.ExecContext(ctx, `
				INSERT INTO sticker_pack (id, name, photo_id)
				VALUES ($1, $2, $3)
			`, packID, packName, stickerID)
			if err != nil {
				return uuid.Nil, err
			}
		} else {
			return uuid.Nil, err
		}
	}

	// Добавление связи sticker — pack
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO sticker_sticker_pack (id, sticker, pack)
		VALUES ($1, $2, $3)
	`, uuid.New(), stickerID, packID)
	if err != nil {
		return uuid.Nil, err
	}

	return stickerID, nil
}

func (r *filesRepository) GetStickerPacks(ctx context.Context) (model.StickerPacks, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, photo_id FROM public.sticker_pack`)
	if err != nil {
		return model.StickerPacks{}, err
	}
	defer rows.Close()

	var packs []model.StickerPack
	for rows.Next() {
		var pack model.StickerPack
		var photoID uuid.UUID

		if err := rows.Scan(&pack.PackID, &pack.Name, &photoID); err != nil {
			log.Println(err)
			return model.StickerPacks{}, err
		}

		// Генерируем URL на превью
		pack.Photo = fmt.Sprintf("/files/%s", photoID.String())
		packs = append(packs, pack)
	}

	if err := rows.Err(); err != nil {
		return model.StickerPacks{}, err
	}
	return model.StickerPacks{Packs: packs}, nil
}

func (r *filesRepository) GetStickerPack(ctx context.Context, packID string) (model.GetStickerPackResponse, error) {
	id, err := uuid.Parse(packID)
	if err != nil {
		return model.GetStickerPackResponse{}, fmt.Errorf("invalid UUID: %w", err)
	}

	var photoPath string
	err = r.db.QueryRowContext(ctx,
		`SELECT photo_id FROM sticker_pack WHERE id = $1`, id).
		Scan(&photoPath)
	if err != nil {
		return model.GetStickerPackResponse{}, err
	}

	// Получаем ID стикеров, входящих в пак, через связь
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id
		FROM sticker_sticker_pack sp
		JOIN sticker s ON s.id = sp.sticker
		WHERE sp.pack = $1
	`, id)
	if err != nil {
		return model.GetStickerPackResponse{}, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var stickerID uuid.UUID
		if err := rows.Scan(&stickerID); err != nil {
			return model.GetStickerPackResponse{}, err
		}
		urls = append(urls, fmt.Sprintf("/files/%s", stickerID.String()))
	}
	if err := rows.Err(); err != nil {
		return model.GetStickerPackResponse{}, err
	}

	return model.GetStickerPackResponse{
		Photo: photoPath, // это уже путь вида "/uploads/stickers/..."
		URLs:  urls,
	}, nil
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
