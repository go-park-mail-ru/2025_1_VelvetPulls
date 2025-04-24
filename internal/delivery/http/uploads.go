package http

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/gorilla/mux"
)

type uploadsController struct {
	fileRepo repository.FileRepository
}

func NewUploadsController(r *mux.Router, fileRepo repository.FileRepository) {
	controller := &uploadsController{
		fileRepo: fileRepo,
	}

	r.HandleFunc("/{folder}/{name}", controller.GetFile).Methods(http.MethodGet)
}

// GetFile отправляет клиенту файл из хранилища (Minio)
//
// @Summary Получение загруженного файла
// @Description Возвращает файл из хранилища
// @Tags Uploads
// @Produce octet-stream
// @Param folder path string true "Название папки"
// @Param name path string true "Имя файла"
// @Success 200 {file} binary
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /uploads/{folder}/{name} [get]
func (c *uploadsController) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	folder := vars["folder"]
	name := vars["name"]
	objectPath := filepath.Join(folder, name)

	// 1. Получаем файл из Minio
	file, err := c.fileRepo.GetFile(ctx, objectPath)
	if err != nil {
		code, appErr := apperrors.GetErrAndCodeToSend(err)
		http.Error(w, appErr.Error(), code)
		return
	}
	defer file.Close()

	// 2. Определяем Content-Type
	contentType := "application/octet-stream"
	if ext := filepath.Ext(name); ext != "" {
		contentType = mime.TypeByExtension(ext)
	}

	// 3. Отправляем файл клиенту
	w.Header().Set("Content-Type", contentType)
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}
