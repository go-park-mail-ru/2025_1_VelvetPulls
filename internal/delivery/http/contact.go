package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type contactController struct {
	sessionUsecase usecase.ISessionUsecase
	contactUsecase usecase.IContactUsecase
}

func NewContactController(r *mux.Router, contactUsecase usecase.IContactUsecase, sessionUsecase usecase.ISessionUsecase) {
	controller := &contactController{
		contactUsecase: contactUsecase,
		sessionUsecase: sessionUsecase,
	}

	r.Handle("/contacts", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetContacts))).Methods(http.MethodGet)
	r.Handle("/contacts", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.AddContact))).Methods(http.MethodPost)
	r.Handle("/contacts", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.DeleteContact))).Methods(http.MethodDelete)
}

// GetContacts получает список контактов пользователя
// @Summary Получить контакты
// @Description Возвращает список контактов пользователя
// @Tags Contacts
// @Produce json
// @Success 200 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /contacts [get]
func (c *contactController) GetContacts(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Fetching contacts")

	contacts, err := c.contactUsecase.GetUserContacts(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get contacts", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusInternalServerError, "Failed to get contacts", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, contacts, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// AddContact добавляет новый контакт
// @Summary Добавить контакт
// @Description Добавляет нового контакта для пользователя
// @Tags Contacts
// @Accept json
// @Produce json
// @Param contact body model.RequestContact true "Данные контакта"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /contacts [post]
func (c *contactController) AddContact(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Adding new contact")

	var contact model.RequestContact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if err := c.contactUsecase.AddUserContact(r.Context(), userID, contact.ID); err != nil {
		logger.Error("Failed to add contact", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusInternalServerError, "Failed to add contact", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, "Contact added successfully", true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// DeleteContact удаляет контакт пользователя
// @Summary Удалить контакт
// @Description Удаляет контакт из списка пользователя
// @Tags Contacts
// @Accept json
// @Produce json
// @Param contact body model.RequestContact true "Данные контакта"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /contacts [delete]
func (c *contactController) DeleteContact(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Deleting contact")

	var contact model.RequestContact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if err := c.contactUsecase.RemoveUserContact(r.Context(), userID, contact.ID); err != nil {
		logger.Error("Failed to delete contact", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusInternalServerError, "Failed to delete contact", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, "Contact deleted successfully", true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}
