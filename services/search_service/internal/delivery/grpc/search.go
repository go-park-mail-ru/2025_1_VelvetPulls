package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/usecase"
	search "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatHandler struct {
	search.UnimplementedChatServiceServer
	chatUC    usecase.ChatUsecase
	contactUC usecase.ContactUsecase
	userUC    usecase.UserUsecase
	messageUC usecase.MessageUsecase
}

func NewChatHandler(
	chatUC usecase.ChatUsecase,
	contactUC usecase.ContactUsecase,
	userUC usecase.UserUsecase,
	messageUC usecase.MessageUsecase,
) *ChatHandler {
	return &ChatHandler{
		chatUC:    chatUC,
		contactUC: contactUC,
		userUC:    userUC,
		messageUC: messageUC,
	}
}

func (h *ChatHandler) SearchUserChats(
	ctx context.Context,
	req *search.SearchUserChatsRequest,
) (*search.SearchUserChatsResponse, error) {
	chatsGrouped, err := h.chatUC.SearchUserChats(ctx, req.UserId, req.Query)
	if err != nil {
		return nil, status.Error(codes.Internal, "search failed")
	}

	resp := &search.SearchUserChatsResponse{
		Dialogs:  make([]*search.Chat, 0),
		Groups:   make([]*search.Chat, 0),
		Channels: make([]*search.Chat, 0),
	}

	if dialogs, ok := chatsGrouped["dialogs"]; ok {
		for _, c := range dialogs {
			resp.Dialogs = append(resp.Dialogs, convertChatToPB(c))
		}
	}

	if groups, ok := chatsGrouped["groups"]; ok {
		for _, c := range groups {
			resp.Groups = append(resp.Groups, convertChatToPB(c))
		}
	}

	if channels, ok := chatsGrouped["channels"]; ok {
		for _, c := range channels {
			resp.Channels = append(resp.Channels, convertChatToPB(c))
		}
	}

	return resp, nil
}

func (h *ChatHandler) SearchContacts(
	ctx context.Context,
	req *search.SearchContactsRequest,
) (*search.SearchContactsResponse, error) {
	const method = "SearchContacts"
	logger := utils.GetLoggerFromCtx(ctx)
	contacts, err := h.contactUC.SearchContacts(ctx, req.GetUserId(), req.GetQuery())
	if err != nil {
		logger.Error(method, zap.Error(err))
		if errors.Is(err, model.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := &search.SearchContactsResponse{
		Contacts: make([]*search.Contact, 0, len(contacts)),
	}

	for _, c := range contacts {
		contactPB := &search.Contact{
			Id:       c.ID.String(),
			Username: c.Username,
		}

		if c.FirstName != nil {
			contactPB.FirstName = *c.FirstName
		}
		if c.LastName != nil {
			contactPB.LastName = *c.LastName
		}
		if c.AvatarURL != nil {
			contactPB.AvatarPath = *c.AvatarURL
		}

		resp.Contacts = append(resp.Contacts, contactPB)
	}
	return resp, nil
}

func (h *ChatHandler) SearchUsers(
	ctx context.Context,
	req *search.SearchUsersRequest,
) (*search.SearchUsersResponse, error) {
	const method = "SearchUsers"
	logger := utils.GetLoggerFromCtx(ctx)
	if len(req.GetQuery()) < 3 {
		return nil, status.Error(codes.InvalidArgument, "search query too short")
	}

	users, err := h.userUC.SearchUsers(ctx, req.GetQuery())
	if err != nil {
		logger.Error(method, zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := &search.SearchUsersResponse{
		Users: make([]*search.User, 0, len(users)),
	}

	for _, u := range users {
		userPB := &search.User{
			Username: u.Username,
		}

		if u.FirstName != nil {
			userPB.FirstName = *u.FirstName
		}
		if u.LastName != nil {
			userPB.LastName = *u.LastName
		}
		if u.AvatarPath != nil {
			userPB.AvatarPath = *u.AvatarPath
		}

		resp.Users = append(resp.Users, userPB)
	}
	return resp, nil
}

func (h *ChatHandler) SearchMessages(
	ctx context.Context,
	req *search.SearchMessagesRequest,
) (*search.SearchMessagesResponse, error) {
	const method = "SearchMessages"
	logger := utils.GetLoggerFromCtx(ctx)
	if req.GetLimit() <= 0 || req.GetLimit() > 100 {
		return nil, status.Error(codes.InvalidArgument, "invalid limit value")
	}

	messages, total, err := h.messageUC.SearchMessages(
		ctx,
		req.GetChatId(),
		req.GetQuery(),
		int(req.GetLimit()),
		int(req.GetOffset()),
	)
	if err != nil {
		logger.Error(method, zap.Error(err))
		if errors.Is(err, model.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, "invalid chat ID format")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := &search.SearchMessagesResponse{
		Messages: make([]*search.Message, 0, len(messages)),
		Total:    int32(total),
	}

	for _, m := range messages {
		msgPB := &search.Message{
			Id:       m.ID.String(),
			Body:     m.Body,
			UserId:   m.UserID.String(),
			SentAt:   m.SentAt.Format(time.RFC3339),
			Username: m.Username,
		}

		resp.Messages = append(resp.Messages, msgPB)
	}
	return resp, nil
}

func convertChatToPB(chat model.Chat) *search.Chat {
	pbChat := &search.Chat{
		Id:           chat.ID.String(),
		Title:        chat.Title,
		Type:         mapChatTypeToProto(chat.Type),
		AvatarPath:   getStringPointer(chat.AvatarPath),
		CreatedAt:    chat.CreatedAt,
		UpdatedAt:    chat.UpdatedAt,
		Participants: convertParticipantsToPB(chat.Participants),
		LastMessage:  convertLastMessageToPB(chat.LastMessage),
	}
	return pbChat
}

// Вспомогательные функции

func mapChatTypeToProto(chatType string) search.ChatType {
	switch chatType {
	case "dialog":
		return search.ChatType_DIALOG
	case "group":
		return search.ChatType_GROUP
	case "channel":
		return search.ChatType_CHANNEL
	default:
		return search.ChatType_DIALOG
	}
}

func mapUserRoleToProto(role string) search.UserRole {
	switch role {
	case "owner":
		return search.UserRole_OWNER
	default:
		return search.UserRole_MEMBER
	}
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func getStringPointer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func convertParticipantsToPB(participants []model.UserInChat) []*search.UserInChat {
	pbParticipants := make([]*search.UserInChat, 0, len(participants))
	for _, p := range participants {
		pbParticipants = append(pbParticipants, &search.UserInChat{
			Id:         p.ID.String(),
			Username:   p.Username,
			AvatarPath: getStringPointer(p.AvatarPath),
		})
	}
	return pbParticipants
}

func convertLastMessageToPB(msg *model.LastMessage) *search.LastMessage {
	if msg == nil {
		return nil
	}

	return &search.LastMessage{
		Id:       msg.ID.String(),
		UserId:   msg.UserID.String(),
		Username: msg.Username,
		Body:     msg.Body,
		SentAt:   msg.SentAt.Format(time.RFC3339),
	}
}
