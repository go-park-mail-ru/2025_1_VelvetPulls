package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/google/uuid"
)

var (
	ErrNotImage    = errors.New("file is not a valid image")
	imageMimeTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
)

func SavePhoto(file multipart.File, folderName string) (string, error) {
	if ok := IsImageFile(file); !ok {
		return "", ErrNotImage
	}

	filenameUUID := uuid.New()

	path := config.UPLOAD_DIR + folderName + "/" + filenameUUID.String() + ".png"
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return path, nil
}

func RewritePhoto(file multipart.File, photoURL string) error {
	dst, err := os.Create(photoURL)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	return nil
}

func RemovePhoto(photoURL string) error {
	err := os.Remove(photoURL)
	if err != nil {
		return err
	}

	return nil
}

func IsImageFile(file multipart.File) bool {
	// Читаем первые 512 байт для определения MIME-типа
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// Определяем MIME-тип
	mimeType := http.DetectContentType(buffer)

	// Сбрасываем позицию чтения
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false
	}

	// Проверяем, что это изображение
	return imageMimeTypes[strings.ToLower(mimeType)]
}
