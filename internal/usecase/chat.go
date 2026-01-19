package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type ChatModule struct {
	userMapper    repository.IUserMapper
	chatMapper    repository.IChatMapper
	messageMapper repository.IMessageMapper
	nc          *nats.Conn
}

type IChatModule interface {
	GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, error)
	GetChatParticipants(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (*model.ChatInfo, error)
	CreateChat(ctx context.Context, userID uuid.UUID, chat *model.CreateChatRequest) (*model.Chat, error)
	UpdateChat(ctx context.Context, userID uuid.UUID, chat *model.UpdateChat) (*model.Chat, error)
	GetChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (*model.Chat, error)
	SendNotifications(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, send bool) error
	DeleteChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error
	AddUsersIntoChat(ctx context.Context, userID uuid.UUID, usernames []string, chatID uuid.UUID) (*model.AddedUsersIntoChat, error)
	SubscribeToChannel(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error
	DeleteUserFromChat(ctx context.Context, userID uuid.UUID, usernamesDelete []string, chatID uuid.UUID) (*model.DeletedUsersFromChat, error)
	LeaveChat(ctx context.Context, userID, chatID uuid.UUID) error
}

func NewChatModule(chatMapper repository.IChatMapper, userMapper repository.IUserMapper, messageMapper repository.IMessageMapper, nc *nats.Conn) IChatModule {
	return &ChatModule{userMapper: userMapper, chatMapper: chatMapper, messageMapper: messageMapper, nc: nc}
}

func (uc *ChatModule) GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetChats", zap.String("userID", userID.String()))

	chats, _, err := uc.chatMapper.GetChats(ctx, userID)
	if err != nil {
		logger.Error("repo.GetChats failed", zap.Error(err))
		return nil, err
	}

	for i := range chats {
		if model.ChatType(chats[i].Type) == model.ChatTypeDialog {
			uc.decorateDialog(ctx, &chats[i], userID)
		}
	}

	logger.Info("GetChats done", zap.Int("count", len(chats)))
	metrics.IncBusinessOp("get_chats")
	return chats, nil
}

func (uc *ChatModule) GetChatParticipants(ctx context.Context, userID, chatID uuid.UUID) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetChat", zap.String("chatID", chatID.String()))

	chat, err := uc.chatMapper.findByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat.Type != string(model.ChatTypeChannel) {
		if err := uc.ensureMember(ctx, userID, chatID); err != nil {
			return nil, err
		}
	}

	role, err := uc.chatMapper.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("GetChatInfo: failed to get user role", zap.Error(err))
		return nil, err
	}

	users, err := uc.userMapper.FindUsersInChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if model.ChatType(chat.Type) == model.ChatTypeDialog {
		uc.decorateDialogInfo(chat, userID, users)
	}

	messages, err := uc.messageMapper.GetMessages(ctx, chatID)
	if err != nil {
		logger.Error("GetChatInfo: failed to get messages", zap.Error(err))
		return nil, err
	}

	metrics.IncBusinessOp("get_Chat")
	return &model.ChatInfo{
		Role:     role,
		Users:    users,
		Messages: messages,
	}, nil
}

func (uc *ChatModule) GetChat(ctx context.Context, userID, chatID uuid.UUID) (*model.Chat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetChat", zap.String("chatID", chatID.String()))

	chat, err := uc.chatMapper.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat.Type != string(model.ChatTypeChannel) {
		if err := uc.ensureMember(ctx, userID, chatID); err != nil {
			return nil, err
		}
	}

	users, err := uc.chatMapper.GetUsersFromChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if model.ChatType(chat.Type) == model.ChatTypeDialog {
		uc.decorateDialogInfo(chat, userID, users)
	}
	sendNotifications, err := uc.chatMapper.GetSendNotifications(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	metrics.IncBusinessOp("get_Chat")
	return &model.Chat{
		ID:                chat.ID,
		AvatarPath:        chat.AvatarPath,
		Type:              chat.Type,
		Title:             chat.Title,
		CountUsers:        len(users),
		SendNotifications: sendNotifications,
	}, nil
}

func (uc *ChatModule) SendNotifications(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, send bool) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("SendNotifications", zap.String("chatID", chatID.String()))
	err := uc.chatMapper.SendNotifications(ctx, userID, chatID, send)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ChatModule) CreateChat(ctx context.Context, userID uuid.UUID, req *model.CreateChatRequest) (*model.Chat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("CreateChat start", zap.String("type", req.Type))

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.Type == string(model.ChatTypeDialog) {
		if info, found := uc.findExistingDialog(ctx, userID, req.Users); found {
			return info, nil
		}
	}

	chatID, err := uc.chatMapper.CreateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	switch model.ChatType(req.Type) {
	case model.ChatTypeDialog:
		if err := uc.addDialogUsers(ctx, userID, req.Users, chatID); err != nil {
			return nil, ErrDialogAddUsers
		}
	case model.ChatTypeGroup:
		if err := uc.addGroupOwner(ctx, userID, req.Users, chatID); err != nil {
			return nil, ErrAddOwnerToGroup
		}
	case model.ChatTypeChannel:
		if err := uc.addChannelOwner(ctx, userID, req.Users, chatID); err != nil {
			return nil, ErrAddOwnerToGroup
		}
	}

	info, err := uc.GetChat(ctx, userID, chatID)

	data, _ := json.Marshal(model.ChatEvent{Action: utils.NewChat, Chat: *info})
	subject := fmt.Sprintf("chat.%s.events", chatID.String())
	if err := uc.nc.Publish(subject, data); err != nil {
		logger.Error("failed to publish chat event", zap.String("subject", subject), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrMessagePublishFailed, err)
	}

	metrics.IncBusinessOp("create_chat")
	return info, err
}

func (uc *ChatModule) UpdateChat(ctx context.Context, userID uuid.UUID, req *model.UpdateChat) (*model.Chat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("UpdateChat", zap.String("chatID", req.ID.String()))

	if req.Avatar != nil && !utils.IsImageFile(*req.Avatar) {
		return nil, utils.ErrNotImage
	}

	chat, err := uc.chatMapper.GetChatByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if model.ChatType(chat.Type) == model.ChatTypeDialog {
		return nil, ErrDialogUpdateForbidden
	}
	if err := uc.ensureOwner(ctx, userID, req.ID); err != nil {
		return nil, err
	}

	newURL, oldURL, err := uc.chatMapper.UpdateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	if req.Avatar != nil && newURL != "" {
		if err := utils.RewritePhoto(*req.Avatar, newURL); err != nil {
			logger.Error("UpdateChat: RewritePhoto failed", zap.Error(err))
			return nil, err
		}
	}

	uc.handleAvatarCleanup(oldURL)

	info, err := uc.GetChat(ctx, userID, req.ID)
	data, _ := json.Marshal(model.ChatEvent{Action: utils.UpdateChat, Chat: *info})

	subject := fmt.Sprintf("chat.%s.events", req.ID.String())
	if err := uc.nc.Publish(subject, data); err != nil {
		logger.Error("failed to publish chat event", zap.String("subject", subject), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrMessagePublishFailed, err)
	}

	metrics.IncBusinessOp("update_chat")
	return info, err
}

func (uc *ChatModule) DeleteChat(ctx context.Context, userID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("DeleteChat", zap.String("chatID", chatID.String()))

	if err := uc.ensureOwner(ctx, userID, chatID); err != nil {
		return err
	}
	if err := uc.chatMapper.DeleteChat(ctx, chatID); err != nil {
		return err
	}

	ce := model.ChatEvent{Action: utils.DeleteChat, Chat: model.Chat{ID: chatID}}
	data, _ := json.Marshal(ce)

	subject := fmt.Sprintf("chat.%s.events", chatID.String())
	if err := uc.nc.Publish(subject, data); err != nil {
		logger.Error("failed to publish chat event", zap.String("subject", subject), zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessagePublishFailed, err)
	}

	metrics.IncBusinessOp("delete_chat")
	return nil
}

func (uc *ChatModule) AddUsersIntoChat(ctx context.Context, userID uuid.UUID, usernames []string, chatID uuid.UUID) (*model.AddedUsersIntoChat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("AddUsersIntoChat", zap.String("chatID", chatID.String()))

	if err := uc.ensureOwner(ctx, userID, chatID); err != nil {
		return nil, err
	}
	chat, err := uc.chatMapper.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if model.ChatType(chat.Type) == model.ChatTypeDialog {
		return nil, ErrDialogAddUsers
	}

	added, notAdded := uc.modifyMembers(ctx, chatID, usernames, true)

	metrics.IncBusinessOp("add_user_into_chat")
	return &model.AddedUsersIntoChat{AddedUsers: added, NotAddedUsers: notAdded}, nil
}

func (uc *ChatModule) SubscribeToChannel(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error {
	chat, err := uc.chatMapper.GetChatByID(ctx, chatID)
	if err != nil {
		return err
	}
	if model.ChatType(chat.Type) != model.ChatTypeChannel {
		return ErrNotChannel
	}
	if err := uc.chatMapper.AddUserToChatByID(ctx, userID, string(model.RoleMember), chatID); err != nil {
		return err
	}
	return nil
}

func (uc *ChatModule) DeleteUserFromChat(ctx context.Context, userID uuid.UUID, usernames []string, chatID uuid.UUID) (*model.DeletedUsersFromChat, error) {
	if err := uc.ensureOwner(ctx, userID, chatID); err != nil {
		return nil, err
	}
	chat, err := uc.chatMapper.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if model.ChatType(chat.Type) == model.ChatTypeDialog {
		return nil, ErrDialogDeleteUsers
	}
	deleted, _ := uc.modifyMembers(ctx, chatID, usernames, false)

	metrics.IncBusinessOp("delete_user_from_chat")
	return &model.DeletedUsersFromChat{DeletedUsers: deleted}, nil
}

func (uc *ChatModule) LeaveChat(ctx context.Context, userID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("LeaveChat", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	role, err := uc.chatMapper.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("LeaveChat: failed to get user role", zap.Error(err))
		return err
	}

	switch model.UserRoleInChat(role) {
	case model.RoleOwner:
		logger.Warn("LeaveChat: owner cannot leave the chat")
		return ErrPermissionDenied

	case model.RoleMember:
		if err := uc.chatMapper.RemoveUserFromChatByID(ctx, userID, chatID); err != nil {
			logger.Error("LeaveChat: failed to remove user from chat", zap.Error(err))
			return err
		}

		event := model.ChatEvent{
			Action: utils.LeaveChat,
			Chat:   model.Chat{ID: chatID},
		}
		data, _ := json.Marshal(event)
		if err := uc.nc.Publish(fmt.Sprintf("chat.%s.events", chatID.String()), data); err != nil {
			logger.Error("LeaveChat: failed to publish event", zap.Error(err))
			return err
		}

		logger.Info("LeaveChat: success")
		metrics.IncBusinessOp("leave_chat")
		return nil

	default:
		logger.Warn("LeaveChat: user is not a member")
		return ErrPermissionDenied
	}
}

// --- Private Helpers ---

func (uc *ChatModule) decorateDialog(ctx context.Context, chat *model.Chat, me uuid.UUID) {
	users, err := uc.chatMapper.GetUsersFromChat(ctx, chat.ID)
	if err != nil {
		zap.L().Warn("decorateDialog: GetUsersFromChat failed", zap.Error(err))
		return
	}

	if len(users) == 1 {
		chat.Title = users[0].Username
		chat.AvatarPath = users[0].AvatarPath
		return
	}

	for _, u := range users {
		if u.ID != me {
			chat.Title = u.Username
			chat.AvatarPath = u.AvatarPath
			break
		}
	}
}

func (uc *ChatModule) decorateDialogInfo(chat *model.Chat, me uuid.UUID, users []model.UserInChat) {

	if len(users) == 1 {
		chat.Title = users[0].Username
		chat.AvatarPath = users[0].AvatarPath
		return
	}

	for _, u := range users {
		if u.ID != me {
			chat.Title = u.Username
			chat.AvatarPath = u.AvatarPath
			break
		}
	}
}

func (uc *ChatModule) ensureMember(ctx context.Context, userID, chatID uuid.UUID) error {
	role, err := uc.chatMapper.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if !model.UserRoleInChat(role).IsMember() {
		return ErrPermissionDenied
	}
	return nil
}

func (uc *ChatModule) ensureOwner(ctx context.Context, userID, chatID uuid.UUID) error {
	role, err := uc.chatMapper.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if model.UserRoleInChat(role) != model.RoleOwner {
		return ErrOnlyOwnerCanModify
	}
	return nil
}

func (uc *ChatModule) findExistingDialog(ctx context.Context, me uuid.UUID, users []string) (*model.Chat, bool) {
	var otherUsername string
	for _, username := range users {
		user, err := uc.userMapper.GetUserByUsername(ctx, username)
		if err != nil {
			return nil, false
		}
		if user.ID != me {
			otherUsername = username
		}
	}
	chats, _, err := uc.chatMapper.GetChats(ctx, me)
	if err != nil {
		zap.L().Warn("findExistingDialog: GetChats failed", zap.Error(err))
		return nil, false
	}
	for _, c := range chats {
		if model.ChatType(c.Type) != model.ChatTypeDialog {
			continue
		}

		users, err := uc.chatMapper.GetUsersFromChat(ctx, c.ID)
		if err != nil {
			continue
		}

		if len(users) == 1 && users[0].ID == me && users[0].Username == otherUsername {
			info, err := uc.GetChat(ctx, me, c.ID)
			if err == nil {
				return info, true
			}
		}

		for _, u := range users {
			if u.ID != me && u.Username == otherUsername {
				info, err := uc.GetChat(ctx, me, c.ID)
				if err == nil {
					return info, true
				}
			}
		}
	}

	return nil, false
}

func (uc *ChatModule) addDialogUsers(ctx context.Context, me uuid.UUID, users []string, chatID uuid.UUID) error {
	err := uc.chatMapper.AddUserToChatByID(ctx, me, string(model.RoleOwner), chatID)
	if err != nil {
		return err
	}
	var otherUsername string
	for _, username := range users {
		user, err := uc.userMapper.GetUserByUsername(ctx, username)
		if err != nil {
			return err
		}
		if user.ID != me {
			otherUsername = username
			err = uc.chatMapper.AddUserToChatByUsername(ctx, otherUsername, string(model.RoleOwner), chatID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (uc *ChatModule) addGroupOwner(ctx context.Context, me uuid.UUID, users []string, chatID uuid.UUID) error {
	err := uc.chatMapper.AddUserToChatByID(ctx, me, string(model.RoleOwner), chatID)
	if err != nil {
		return err
	}

	var otherUsername string
	for _, username := range users {
		user, err := uc.userMapper.GetUserByUsername(ctx, username)
		if err != nil {
			return err
		}
		if user.ID != me {
			otherUsername = username
			err = uc.chatMapper.AddUserToChatByUsername(ctx, otherUsername, string(model.RoleMember), chatID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *ChatModule) addChannelOwner(ctx context.Context, ownerID uuid.UUID, users []string, chatID uuid.UUID) error {
	err := uc.chatMapper.AddUserToChatByID(ctx, ownerID, string(model.RoleOwner), chatID)
	if err != nil {
		return err
	}
	var otherUsername string
	for _, username := range users {
		user, err := uc.userMapper.GetUserByUsername(ctx, username)
		if err != nil {
			return err
		}
		if user.ID != ownerID {
			otherUsername = username
			err = uc.chatMapper.AddUserToChatByUsername(ctx, otherUsername, string(model.RoleMember), chatID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (uc *ChatModule) modifyMembers(ctx context.Context, chatID uuid.UUID, names []string, add bool) ([]string, []string) {
	var (
		success []string
		failed  []string
	)
	for _, name := range names {
		var err error
		if add {
			err = uc.chatMapper.AddUserToChatByUsername(ctx, name, string(model.RoleMember), chatID)
		} else {
			err = uc.chatMapper.RemoveUserFromChatByUsername(ctx, name, chatID)
		}
		if err != nil {
			failed = append(failed, name)
		} else {
			success = append(success, name)
		}
	}
	return success, failed
}

// handleAvatarCleanup removes only old avatar
func (uc *ChatModule) handleAvatarCleanup(oldURL string) {
	if oldURL != "" {
		go func(url string) {
			if err := utils.RemovePhoto(url); err != nil {
				zap.L().Warn("Old avatar remove failed", zap.String("url", url), zap.Error(err))
			}
		}(oldURL)
	}
}
