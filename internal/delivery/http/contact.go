package http

import (
	"net/http"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type contactController struct {
	sessionClient  authpb.SessionServiceClient
	contactUsecase usecase.IContactUsecase
}

func NewContactController(r *mux.Router, contactUsecase usecase.IContactUsecase, sessionClient authpb.SessionServiceClient) {
	controller := &contactController{
		contactUsecase: contactUsecase,
		sessionClient:  sessionClient,
	}

	r.Handle("/contacts", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetContacts))).Methods(http.MethodGet)
	r.Handle("/contacts", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.AddContact))).Methods(http.MethodPost)
	r.Handle("/contacts", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.DeleteContact))).Methods(http.MethodDelete)
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

	logger.Info("GetContacts called", zap.String("userID", userID.String()))

	contacts, err := c.contactUsecase.GetUserContacts(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get contacts", zap.Error(err))
		code, errToSend := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, errToSend, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, contacts, true)
}

// AddContact добавляет новый контакт
// @Summary Добавить контакт
// @Description Добавляет нового контакта для пользователя
// @Tags Contacts
// @Accept json
// @Produce json
// @Param contact body model.RequestContact true "Данные контакта"
// @Success 200 {object} model.Contact
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /contacts [post]
func (c *contactController) AddContact(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())

	logger.Info("AddContact called", zap.String("userID", userID.String()))

	var contact model.RequestContact
	if err := utils.ParseJSONRequest(r, &contact); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request body", false)
		return
	}

	addedContact, err := c.contactUsecase.AddUserContact(r.Context(), userID, contact.Username)
	if err != nil {
		logger.Error("Failed to add contact", zap.Error(err))
		code, errToSend := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, errToSend, false)
		return
	}

	// Возвращаем добавленного пользователя
	addedContact.Sanitize()
	utils.SendJSONResponse(w, r, http.StatusOK, addedContact, true)
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

	logger.Info("DeleteContact called", zap.String("userID", userID.String()))

	var contact model.RequestContact
	if err := utils.ParseJSONRequest(r, &contact); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request body", false)
		return
	}

	if err := c.contactUsecase.RemoveUserContact(r.Context(), userID, contact.Username); err != nil {
		logger.Error("Failed to delete contact", zap.Error(err))
		code, errToSend := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, errToSend, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, "Contact deleted successfully", true)
}
