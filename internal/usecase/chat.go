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

type ChatUsecase struct {
	chatRepo repository.IChatRepo
	nc       *nats.Conn
}

type IChatUsecase interface {
	GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, error)
	GetChatInfo(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (*model.ChatInfo, error)
	CreateChat(ctx context.Context, userID uuid.UUID, chat *model.CreateChat) (*model.ChatInfo, error)
	UpdateChat(ctx context.Context, userID uuid.UUID, chat *model.UpdateChat) (*model.ChatInfo, error)
	DeleteChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error
	AddUsersIntoChat(ctx context.Context, userID uuid.UUID, usernames []string, chatID uuid.UUID) (*model.AddedUsersIntoChat, error)
	DeleteUserFromChat(ctx context.Context, userID uuid.UUID, usernamesDelete []string, chatID uuid.UUID) (*model.DeletedUsersFromChat, error)
	LeaveChat(ctx context.Context, userID, chatID uuid.UUID) error
}

func NewChatUsecase(chatRepo repository.IChatRepo, nc *nats.Conn) IChatUsecase {
	return &ChatUsecase{chatRepo: chatRepo, nc: nc}
}

func (uc *ChatUsecase) GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetChats", zap.String("userID", userID.String()))

	chats, _, err := uc.chatRepo.GetChats(ctx, userID)
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

func (uc *ChatUsecase) GetChatInfo(ctx context.Context, userID, chatID uuid.UUID) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetChatInfo", zap.String("chatID", chatID.String()))

	if err := uc.ensureMember(ctx, userID, chatID); err != nil {
		return nil, err
	}

	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	users, err := uc.chatRepo.GetUsersFromChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if model.ChatType(chat.Type) == model.ChatTypeDialog {
		uc.decorateDialogInfo(chat, userID, users)
	}
	metrics.IncBusinessOp("get_chatinfo")
	return &model.ChatInfo{
		ID:         chat.ID,
		AvatarPath: chat.AvatarPath,
		Type:       chat.Type,
		Title:      chat.Title,
		CountUsers: len(users),
		Users:      users,
	}, nil
}

func (uc *ChatUsecase) CreateChat(ctx context.Context, userID uuid.UUID, req *model.CreateChat) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("CreateChat start", zap.String("type", req.Type))

	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Avatar != nil && !utils.IsImageFile(*req.Avatar) {
		return nil, utils.ErrNotImage
	}

	if req.Type == string(model.ChatTypeDialog) {
		if info, found := uc.findExistingDialog(ctx, userID, req.DialogUser); found {
			return info, nil
		}
	}

	chatID, avatarURL, err := uc.chatRepo.CreateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	if req.Avatar != nil {
		if err := utils.RewritePhoto(*req.Avatar, avatarURL); err != nil {
			logger.Error("CreateChat: RewritePhoto failed", zap.Error(err))
			return nil, err
		}
	}

	switch model.ChatType(req.Type) {
	case model.ChatTypeDialog:
		if err := uc.addDialogUsers(ctx, userID, req.DialogUser, chatID); err != nil {
			return nil, ErrDialogAddUsers
		}
	case model.ChatTypeGroup:
		if err := uc.addGroupOwner(ctx, userID, chatID); err != nil {
			return nil, ErrAddOwnerToGroup
		}
	case model.ChatTypeChannel:
		if err := uc.addChannelOwner(ctx, userID, chatID); err != nil {
			return nil, ErrAddOwnerToGroup
		}
	}

	info, err := uc.GetChatInfo(ctx, userID, chatID)

	data, _ := json.Marshal(model.ChatEvent{Action: utils.NewChat, Chat: *info})
	uc.nc.Publish(fmt.Sprintf("chat.%s.events", chatID.String()), data)

	metrics.IncBusinessOp("create_chat")
	return info, err
}

func (uc *ChatUsecase) UpdateChat(ctx context.Context, userID uuid.UUID, req *model.UpdateChat) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("UpdateChat", zap.String("chatID", req.ID.String()))

	if req.Avatar != nil && !utils.IsImageFile(*req.Avatar) {
		return nil, utils.ErrNotImage
	}

	chat, err := uc.chatRepo.GetChatByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if model.ChatType(chat.Type) == model.ChatTypeDialog {
		return nil, ErrDialogUpdateForbidden
	}
	if err := uc.ensureOwner(ctx, userID, req.ID); err != nil {
		return nil, err
	}

	newURL, oldURL, err := uc.chatRepo.UpdateChat(ctx, req)
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

	info, err := uc.GetChatInfo(ctx, userID, req.ID)
	data, _ := json.Marshal(model.ChatEvent{Action: utils.UpdateChat, Chat: *info})
	uc.nc.Publish(fmt.Sprintf("chat.%s.events", req.ID.String()), data)

	metrics.IncBusinessOp("update_chat")
	return info, err
}

func (uc *ChatUsecase) DeleteChat(ctx context.Context, userID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("DeleteChat", zap.String("chatID", chatID.String()))

	if err := uc.ensureOwner(ctx, userID, chatID); err != nil {
		return err
	}
	if err := uc.chatRepo.DeleteChat(ctx, chatID); err != nil {
		return err
	}

	ce := model.ChatEvent{Action: utils.DeleteChat, Chat: model.ChatInfo{ID: chatID}}
	data, _ := json.Marshal(ce)
	uc.nc.Publish(fmt.Sprintf("chat.%s.events", chatID.String()), data)
	metrics.IncBusinessOp("delete_chat")
	return nil
}

func (uc *ChatUsecase) AddUsersIntoChat(ctx context.Context, userID uuid.UUID, usernames []string, chatID uuid.UUID) (*model.AddedUsersIntoChat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("AddUsersIntoChat", zap.String("chatID", chatID.String()))

	if err := uc.ensureOwner(ctx, userID, chatID); err != nil {
		return nil, err
	}
	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
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

func (uc *ChatUsecase) DeleteUserFromChat(ctx context.Context, userID uuid.UUID, usernames []string, chatID uuid.UUID) (*model.DeletedUsersFromChat, error) {
	if err := uc.ensureOwner(ctx, userID, chatID); err != nil {
		return nil, err
	}
	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
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

func (uc *ChatUsecase) LeaveChat(ctx context.Context, userID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("LeaveChat", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("LeaveChat: failed to get user role", zap.Error(err))
		return err
	}

	switch model.UserRoleInChat(role) {
	case model.RoleOwner:
		logger.Warn("LeaveChat: owner cannot leave the chat")
		return ErrPermissionDenied

	case model.RoleMember:
		if err := uc.chatRepo.RemoveUserFromChatByID(ctx, userID, chatID); err != nil {
			logger.Error("LeaveChat: failed to remove user from chat", zap.Error(err))
			return err
		}

		event := model.ChatEvent{
			Action: utils.LeaveChat,
			Chat:   model.ChatInfo{ID: chatID},
		}
		data, _ := json.Marshal(event)
		uc.nc.Publish(fmt.Sprintf("chat.%s.events", chatID.String()), data)

		logger.Info("LeaveChat: success")
		metrics.IncBusinessOp("leave_chat")
		return nil

	default:
		logger.Warn("LeaveChat: user is not a member")
		return ErrPermissionDenied
	}
}

// --- Private Helpers ---

func (uc *ChatUsecase) decorateDialog(ctx context.Context, chat *model.Chat, me uuid.UUID) {
	users, err := uc.chatRepo.GetUsersFromChat(ctx, chat.ID)
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

func (uc *ChatUsecase) decorateDialogInfo(chat *model.Chat, me uuid.UUID, users []model.UserInChat) {

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

func (uc *ChatUsecase) ensureMember(ctx context.Context, userID, chatID uuid.UUID) error {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if !model.UserRoleInChat(role).IsMember() {
		return ErrPermissionDenied
	}
	return nil
}

func (uc *ChatUsecase) ensureOwner(ctx context.Context, userID, chatID uuid.UUID) error {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if model.UserRoleInChat(role) != model.RoleOwner {
		return ErrOnlyOwnerCanModify
	}
	return nil
}

func (uc *ChatUsecase) findExistingDialog(ctx context.Context, me uuid.UUID, otherUsername string) (*model.ChatInfo, bool) {
	chats, _, err := uc.chatRepo.GetChats(ctx, me)
	if err != nil {
		zap.L().Warn("findExistingDialog: GetChats failed", zap.Error(err))
		return nil, false
	}
	for _, c := range chats {
		if model.ChatType(c.Type) != model.ChatTypeDialog {
			continue
		}

		users, err := uc.chatRepo.GetUsersFromChat(ctx, c.ID)
		if err != nil {
			continue
		}

		if len(users) == 1 && users[0].ID == me && users[0].Username == otherUsername {
			info, err := uc.GetChatInfo(ctx, me, c.ID)
			if err == nil {
				return info, true
			}
		}

		for _, u := range users {
			if u.ID != me && u.Username == otherUsername {
				info, err := uc.GetChatInfo(ctx, me, c.ID)
				if err == nil {
					return info, true
				}
			}
		}
	}

	return nil, false
}

func (uc *ChatUsecase) addDialogUsers(ctx context.Context, me uuid.UUID, otherUsername string, chatID uuid.UUID) error {
	err := uc.chatRepo.AddUserToChatByID(ctx, me, string(model.RoleOwner), chatID)
	if err != nil {
		return err
	}

	err = uc.chatRepo.AddUserToChatByUsername(ctx, otherUsername, string(model.RoleMember), chatID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ChatUsecase) addGroupOwner(ctx context.Context, me uuid.UUID, chatID uuid.UUID) error {
	err := uc.chatRepo.AddUserToChatByID(ctx, me, string(model.RoleOwner), chatID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ChatUsecase) addChannelOwner(ctx context.Context, ownerID uuid.UUID, chatID uuid.UUID) error {
	err := uc.chatRepo.AddUserToChatByID(ctx, ownerID, string(model.RoleOwner), chatID)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ChatUsecase) modifyMembers(ctx context.Context, chatID uuid.UUID, names []string, add bool) ([]string, []string) {
	var (
		success []string
		failed  []string
	)
	for _, name := range names {
		var err error
		if add {
			err = uc.chatRepo.AddUserToChatByUsername(ctx, name, string(model.RoleMember), chatID)
		} else {
			err = uc.chatRepo.RemoveUserFromChatByUsername(ctx, name, chatID)
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
func (uc *ChatUsecase) handleAvatarCleanup(oldURL string) {
	if oldURL != "" {
		go func(url string) {
			if err := utils.RemovePhoto(url); err != nil {
				zap.L().Warn("Old avatar remove failed", zap.String("url", url), zap.Error(err))
			}
		}(oldURL)
	}
}
