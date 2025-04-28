package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/google/uuid"
)

var (
	ErrNotImage      = errors.New("file is not a valid image")
	ErrSavingImage   = errors.New("failed to save image")
	ErrUpdatingImage = errors.New("failed to update image")
	ErrDeletingImage = errors.New("failed to delete image")
	imageMimeTypes   = map[string]bool{
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
		return "", ErrSavingImage
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", ErrSavingImage
	}

	return path, nil
}

func RewritePhoto(file multipart.File, photoURL string) error {
	dir := filepath.Dir(photoURL)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	dst, err := os.Create(photoURL)
	if err != nil {
		return ErrUpdatingImage
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return ErrUpdatingImage
	}

	return nil
}

func RemovePhoto(photoURL string) error {
	if err := os.Remove(photoURL); err != nil {
		return ErrDeletingImage
	}
	return nil
}

func IsImageFile(file multipart.File) bool {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil && err != io.EOF {
		return false
	}

	mimeType := http.DetectContentType(buffer)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false
	}

	return imageMimeTypes[strings.ToLower(mimeType)]
}
