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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatHandler struct {
	search.UnimplementedChatServiceServer
	chatUC    usecase.ChatUsecase
	contactUC usecase.ContactUsecase
	messageUC usecase.MessageUsecase
}

func NewChatHandler(
	chatUC usecase.ChatUsecase,
	contactUC usecase.ContactUsecase,
	messageUC usecase.MessageUsecase,
) *ChatHandler {
	return &ChatHandler{
		chatUC:    chatUC,
		contactUC: contactUC,
		messageUC: messageUC,
	}
}

func (h *ChatHandler) SearchUserChats(
	ctx context.Context,
	req *search.SearchUserChatsRequest,
) (*search.SearchUserChatsResponse, error) {
	const method = "SearchUserChats"
	logger := utils.GetLoggerFromCtx(ctx)
	globalChannels, groups, err := h.chatUC.SearchUserChats(ctx, req.GetUserId(), req.GetQuery())
	if err != nil {
		logger.Error(method, zap.Error(err))
		if errors.Is(err, model.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := &search.SearchUserChatsResponse{
		GlobalChannels: convertChatsToPB(globalChannels),
		Dialogs:        convertChatsToPB(groups.Dialogs),
		Groups:         convertChatsToPB(groups.Groups),
		Channels:       convertChatsToPB(groups.Channels),
	}
	return resp, nil
}

func (h *ChatHandler) SearchContacts(
	ctx context.Context,
	req *search.SearchContactsRequest,
) (*search.SearchContactsResponse, error) {
	const method = "SearchContacts"
	logger := utils.GetLoggerFromCtx(ctx)
	users, contacts, err := h.contactUC.SearchContacts(ctx, req.GetUserId(), req.GetQuery())
	if err != nil {
		logger.Error(method, zap.Error(err))
		if errors.Is(err, model.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	resp := &search.SearchContactsResponse{
		Contacts: make([]*search.Contact, 0, len(contacts)),
		Users:    make([]*search.User, 0, len(users)),
	}

	for _, c := range contacts {
		contactPB := &search.Contact{
			Id:       c.ID.String(),
			Username: c.Username,
		}

		if c.Name != nil {
			contactPB.Name = *c.Name
		}
		if c.AvatarURL != nil {
			contactPB.AvatarPath = *c.AvatarURL
		}

		resp.Contacts = append(resp.Contacts, contactPB)
	}

	for _, u := range users {
		userPB := &search.User{
			Id:       u.ID.String(),
			Username: u.Username,
		}

		if u.Name != nil {
			userPB.Name = *u.Name
		}
		if u.AvatarPath != nil {
			userPB.AvatarPath = *u.AvatarPath
		}

		if u.BirthDate != nil {
			userPB.BirthDate = timestamppb.New(*u.BirthDate)
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

func convertChatsToPB(chats []model.Chat) []*search.Chat {
	pbChats := make([]*search.Chat, 0, len(chats))
	for _, c := range chats {
		chatPB := &search.Chat{
			Id:        c.ID.String(),
			Title:     c.Title,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}

		switch c.Type {
		case "dialog":
			chatPB.Type = search.ChatType_DIALOG
		case "group":
			chatPB.Type = search.ChatType_GROUP
		case "channel":
			chatPB.Type = search.ChatType_CHANNEL
		}

		if c.AvatarPath != nil {
			chatPB.AvatarPath = *c.AvatarPath
		}

		if c.LastMessage != nil {
			lastMsgPB := &search.LastMessage{
				Id:       c.LastMessage.ID.String(),
				UserId:   c.LastMessage.UserID.String(),
				Body:     c.LastMessage.Body,
				SentAt:   c.LastMessage.SentAt.Format(time.RFC3339),
				Username: c.LastMessage.Username,
			}
			chatPB.LastMessage = lastMsgPB
		}

		pbChats = append(pbChats, chatPB)
	}
	return pbChats
}
