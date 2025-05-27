package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type filesController struct {
	sessionClient authpb.SessionServiceClient
	filesUsecase  usecase.IFilesUsecase
}

func NewFilesController(r *mux.Router, sessionClient authpb.SessionServiceClient, filesUsecase usecase.IFilesUsecase) {
	URLs := []string{
		"/uploads/stickers/gopher/6762d5505803e3d181d0ecc9.webp",
		"/uploads/stickers/gopher/6762d7f95803e3d181d0ecca.webp",
		"/uploads/stickers/gopher/6762d8aa5803e3d181d0eccb.webp",
		"/uploads/stickers/gopher/6762d8d85803e3d181d0eccc.webp",
		"/uploads/stickers/gopher/6762d8f45803e3d181d0eccd.webp",
		"/uploads/stickers/gopher/6762d90e5803e3d181d0ecce.webp",
		"/uploads/stickers/gopher/6762d9215803e3d181d0eccf.webp",
	}

	for _, url := range URLs {
		file, header, err := getMultipartFile(url)
		if err != nil {
			log.Printf("stickers get error: %v", err)
		}

		err = filesUsecase.SaveSticker(context.Background(), file, header, "gopher")
		if err != nil {
			log.Printf("stickers save error: %v", err)
		}
	}

	controller := &filesController{
		sessionClient: sessionClient,
		filesUsecase:  filesUsecase,
	}

	r.Handle("/files/{file_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetFile))).Methods(http.MethodGet)
	r.Handle("/stickerpacks/{pack_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetStickerPack))).Methods(http.MethodGet)
	r.Handle("/stickerpacks", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetStickerPacks))).Methods(http.MethodGet)
}

type File struct {
	*os.File
}

func (f *File) Close() error {
	return f.File.Close()
}

func newFileHeader(filePath string) (*multipart.FileHeader, error) {
	// Получаем информацию о файле
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Заполняем заголовок
	header := &multipart.FileHeader{
		Filename: fileInfo.Name(),
		Size:     fileInfo.Size(),
		Header:   make(textproto.MIMEHeader),
	}
	header.Header.Set("Content-Type", "image/webp")
	return header, nil
}

func getMultipartFile(filePath string) (multipart.File, *multipart.FileHeader, error) {
	// Открываем файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("файл не существует: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close() // Закрываем файл при ошибке
		return nil, nil, err
	}
	if fileInfo.Size() == 0 {
		file.Close() // Закрываем файл при отсутствии данных
		return nil, nil, fmt.Errorf("файл пустой: %s", filePath)
	}

	// Получаем заголовок
	header, err := newFileHeader(filePath)
	if err != nil {
		file.Close() // Закрываем файл при ошибке
		return nil, nil, err
	}

	return &File{File: file}, header, nil
}

// GetFile godoc
// @Summary Получить файл
// @Description Получить файл по его ID
// @Tags files
// @Accept json
// @Produce octet-stream
// @Param file_id path string true "File ID"
// @Success 200 {file} file "Файл успешно получен"
// @Failure 404 {object} utils.ErrorResponse "Файл не найден"
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера"
// @Router /files/{file_id} [get]
func (c *filesController) GetFile(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	fileID, parseErr := uuid.Parse(mux.Vars(r)["file_id"])
	if parseErr != nil {
		logger.Error("Invalid chat ID format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	file, fileData, err := c.filesUsecase.GetFile(r.Context(), fileID, userID)
	if err != nil {
		logger.Error("Failed to get file",
			zap.String("fileID", fileID.String()),
			zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusNotFound, "File not found", false)
		return
	}

	w.Header().Set("Content-Type", fileData.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileData.Filename))

	if _, err := io.Copy(w, file); err != nil {
		logger.Error("Failed to send file",
			zap.String("fileID", fileID.String()),
			zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusInternalServerError, "Failed to send file", false)
	}
}

// GetStickerPack godoc
// @Summary Получить стикерпак
// @Description Получить стикерпак по его ID
// @Tags stickers
// @Accept json
// @Produce json
// @Param pack_id path string true "ID стикерпака"
// @Success 200 {object} model.GetStickerPackResponse
// @Failure 404 {object} utils.ErrorResponse "Стикерпак не найден"
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера"
// @Router /stickerpacks/{pack_id} [get]
func (c *filesController) GetStickerPack(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	vars := mux.Vars(r)
	packID := vars["pack_id"]

	response, err := c.filesUsecase.GetStickerPack(r.Context(), packID)
	if err != nil {
		logger.Error("Failed to get sticker pack",
			zap.String("packID", packID),
			zap.Error(err))

		status := http.StatusInternalServerError
		if err.Error() == "sticker pack not found" {
			status = http.StatusNotFound
		}

		utils.SendJSONResponse(w, r, status, err.Error(), false)
		return
	}
	resp, err := easyjson.Marshal(response)
	if err != nil {
		logger.Error("Marshaling error", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusInternalServerError, "Internal error", false)
		return
	}
	utils.SendJSONResponse(w, r, http.StatusOK, resp, true)
}

// GetStickerPacks godoc
// @Summary Получить все стикерпаки
// @Description Получить список всех доступных стикерпаков
// @Tags stickers
// @Accept json
// @Produce json
// @Success 200 {object} model.StickerPacks
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера"
// @Router /stickerpacks [get]
func (c *filesController) GetStickerPacks(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	packs, err := c.filesUsecase.GetStickerPacks(r.Context())
	if err != nil {
		logger.Error("Failed to get sticker packs",
			zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusInternalServerError, "Failed to get sticker packs", false)
		return
	}
	resp, err := easyjson.Marshal(packs)
	if err != nil {
		logger.Error("Marshaling error", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusInternalServerError, "Internal error", false)
		return
	}
	utils.SendJSONResponse(w, r, http.StatusOK, resp, true)
}
