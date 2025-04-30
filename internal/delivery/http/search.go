package http

import (
	"net/http"
	"strconv"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/delivery/proto"
	chatpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/delivery/proto"
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
	r.Handle("/search", http.HandlerFunc(controller.SearchChats)).Methods(http.MethodGet)
	r.Handle("/search/{chat_id}/messages", http.HandlerFunc(controller.SearchMessages)).Methods(http.MethodGet)

	// Контакты
	r.Handle("/search/contacts", http.HandlerFunc(controller.SearchContacts)).Methods(http.MethodGet)
}

// func (c *searchController) AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		sessionID, err := utils.GetSessionCookie(r)
// 		if err != nil {
// 			utils.SendJSONResponse(w, r, http.StatusUnauthorized, "Unauthorized", false)
// 			return
// 		}

// 		userID, err := c.sessionClient.GetSession(r.Context(), &authpb.GetSessionRequest{
// 			SessionId: sessionID,
// 		})
// 		if err != nil {
// 			code, msg := apperrors.UnpackGrpcError(err)
// 			utils.SendJSONResponse(w, r, code, msg, false)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "userID", userID.GetUserId())
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func (c *searchController) SearchChats(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := r.Context().Value("userID").(string)

	query := r.URL.Query().Get("query")
	types := r.URL.Query()["type"]

	resp, err := c.searchClient.SearchUserChats(r.Context(), &chatpb.SearchUserChatsRequest{
		UserId: userID,
		Query:  query,
		Types:  types,
	})
	if err != nil {
		logger.Error("gRPC SearchChats error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, resp.Chats, true)
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
	userID := r.Context().Value("userID").(string)

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

	utils.SendJSONResponse(w, r, http.StatusOK, resp.Contacts, true)
}
