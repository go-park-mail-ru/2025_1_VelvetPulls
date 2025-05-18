package http

import (
	"net/http"
	"strconv"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	mw "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	chatpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/proto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type searchController struct {
	searchClient  chatpb.ChatServiceClient
	sessionClient authpb.SessionServiceClient
}

func NewSearchController(r *mux.Router, searchClient chatpb.ChatServiceClient, sessionClient authpb.SessionServiceClient) {
	controller := &searchController{
		searchClient:  searchClient,
		sessionClient: sessionClient,
	}

	// Чат
	r.Handle("/search", mw.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.SearchChats))).Methods(http.MethodGet)
	r.Handle("/search/{chat_id}/messages", mw.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.SearchMessages))).Methods(http.MethodGet)

	// Контакты
	r.Handle("/search/contacts", mw.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.SearchContacts))).Methods(http.MethodGet)
}

func (c *searchController) SearchChats(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context()).String()

	query := r.URL.Query().Get("query")

	resp, err := c.searchClient.SearchUserChats(r.Context(), &chatpb.SearchUserChatsRequest{
		UserId: userID,
		Query:  query,
	})
	if err != nil {
		logger.Error("gRPC SearchChats error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, resp, true)
}

func (c *searchController) SearchMessages(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	vars := mux.Vars(r)
	chatID := vars["chat_id"]

	query := r.URL.Query().Get("query")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	resp, err := c.searchClient.SearchMessages(r.Context(), &chatpb.SearchMessagesRequest{
		ChatId: chatID,
		Query:  query,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		logger.Error("gRPC SearchMessages error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, map[string]interface{}{
		"messages": resp.Messages,
		"total":    resp.Total,
	}, true)
}

func (c *searchController) SearchContacts(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context()).String()

	query := r.URL.Query().Get("query")

	resp, err := c.searchClient.SearchContacts(r.Context(), &chatpb.SearchContactsRequest{
		UserId: userID,
		Query:  query,
	})
	if err != nil {
		logger.Error("gRPC SearchContacts error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, resp, true)
}
