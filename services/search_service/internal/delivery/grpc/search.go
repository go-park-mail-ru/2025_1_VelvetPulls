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
	const method = "SearchUserChats"
	logger := utils.GetLoggerFromCtx(ctx)
	chats, err := h.chatUC.SearchUserChats(ctx, req.GetUserId(), req.GetQuery(), req.GetTypes())
	if err != nil {
		logger.Error(method, zap.Error(err))
		if errors.Is(err, model.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := &search.SearchUserChatsResponse{
		Chats: make([]*search.Chat, 0, len(chats)),
	}

	for _, c := range chats {
		chatPB := &search.Chat{
			Id:        c.ID.String(),
			Type:      c.Type,
			Title:     c.Title,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}

		if c.AvatarPath != nil {
			chatPB.AvatarPath = *c.AvatarPath
		}

		if c.LastMessage != nil {
			chatPB.LastMessage = &search.LastMessage{
				Id:       c.LastMessage.ID.String(),
				UserId:   c.LastMessage.UserID.String(),
				Body:     c.LastMessage.Body,
				SentAt:   c.LastMessage.SentAt.Format(time.RFC3339),
				Username: c.LastMessage.Username,
			}
		}

		resp.Chats = append(resp.Chats, chatPB)
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
