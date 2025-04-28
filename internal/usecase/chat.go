package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChatUsecase struct {
	chatRepo  repository.IChatRepo
	wsUsecase IWebsocketUsecase
}

// IChatUsecase describes chat operations
type IChatUsecase interface {
	GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, error)
	GetChatInfo(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (*model.ChatInfo, error)
	CreateChat(ctx context.Context, userID uuid.UUID, chat *model.CreateChat) (*model.ChatInfo, error)
	UpdateChat(ctx context.Context, userID uuid.UUID, chat *model.UpdateChat) (*model.ChatInfo, error)
	DeleteChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error
	AddUsersIntoChat(ctx context.Context, userID uuid.UUID, usernames []string, chatID uuid.UUID) (*model.AddedUsersIntoChat, error)
	DeleteUserFromChat(ctx context.Context, userID uuid.UUID, usernamesDelete []string, chatID uuid.UUID) (*model.DeletedUsersFromChat, error)
}

// NewChatUsecase constructs ChatUsecase
func NewChatUsecase(chatRepo repository.IChatRepo, wsUsecase IWebsocketUsecase) IChatUsecase {
	return &ChatUsecase{chatRepo: chatRepo, wsUsecase: wsUsecase}
}

// GetChats returns summary for each chat
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
	return chats, nil
}

// GetChatInfo loads full info and enforces membership
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

	return &model.ChatInfo{
		ID:         chat.ID,
		AvatarPath: chat.AvatarPath,
		Type:       chat.Type,
		Title:      chat.Title,
		CountUsers: len(users),
		Users:      users,
	}, nil
}

// CreateChat creates new chat, adds users, initializes WS and publishes event
func (uc *ChatUsecase) CreateChat(ctx context.Context, userID uuid.UUID, req *model.CreateChat) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("CreateChat start", zap.String("type", req.Type))

	// Validation
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Avatar != nil && !utils.IsImageFile(*req.Avatar) {
		return nil, utils.ErrNotImage
	}

	// For dialogs, reuse existing
	if req.Type == string(model.ChatTypeDialog) {
		if info, found := uc.findExistingDialog(ctx, userID, req.DialogUser); found {
			return info, nil
		}
	}

	// Create chat record
	chatID, avatarURL, err := uc.chatRepo.CreateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	// Save avatar file if provided
	if req.Avatar != nil {
		if err := utils.RewritePhoto(*req.Avatar, avatarURL); err != nil {
			logger.Error("CreateChat: RewritePhoto failed", zap.Error(err))
			return nil, err
		}
	}

	// Add initial participant(s)
	switch model.ChatType(req.Type) {
	case model.ChatTypeDialog:
		uc.addDialogUsers(ctx, userID, req.DialogUser, chatID)
	case model.ChatTypeGroup:
		uc.addGroupOwner(ctx, userID, chatID)
	case model.ChatTypeChannel:
		uc.addChannelOwner(ctx, userID, chatID)
	}

	// Prepare ChatInfo and publish event
	info, err := uc.GetChatInfo(ctx, userID, chatID)
	if err == nil && uc.wsUsecase != nil {
		uc.wsUsecase.InitChat(chatID, uc.extractUserIDs(info.Users))
		uc.wsUsecase.PublishChatEvent(model.ChatEvent{Action: NewChat, Chat: *info})
	}

	return info, err
}

// UpdateChat updates chat metadata
func (uc *ChatUsecase) UpdateChat(ctx context.Context, userID uuid.UUID, req *model.UpdateChat) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("UpdateChat", zap.String("chatID", req.ID.String()))

	if req.Avatar != nil && !utils.IsImageFile(*req.Avatar) {
		return nil, utils.ErrNotImage
	}

	// ownership & existence
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

	// update DB
	newURL, oldURL, err := uc.chatRepo.UpdateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	// save new avatar file if present
	if req.Avatar != nil && newURL != "" {
		if err := utils.RewritePhoto(*req.Avatar, newURL); err != nil {
			logger.Error("UpdateChat: RewritePhoto failed", zap.Error(err))
			return nil, err
		}
	}

	// cleanup old file
	uc.handleAvatarCleanup(oldURL)

	// publish update event
	info, err := uc.GetChatInfo(ctx, userID, req.ID)
	if err == nil && uc.wsUsecase != nil {
		uc.wsUsecase.PublishChatEvent(model.ChatEvent{Action: UpdateChat, Chat: *info})
	}

	return info, err
}

// DeleteChat removes chat
func (uc *ChatUsecase) DeleteChat(ctx context.Context, userID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("DeleteChat", zap.String("chatID", chatID.String()))

	if err := uc.ensureOwner(ctx, userID, chatID); err != nil {
		return err
	}
	if err := uc.chatRepo.DeleteChat(ctx, chatID); err != nil {
		return err
	}
	if uc.wsUsecase != nil {
		uc.wsUsecase.PublishChatEvent(model.ChatEvent{Action: DeleteChat, Chat: model.ChatInfo{ID: chatID}})
	}
	return nil
}

// AddUsersIntoChat adds members
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

	return &model.AddedUsersIntoChat{AddedUsers: added, NotAddedUsers: notAdded}, nil
}

// DeleteUserFromChat removes members
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
	return &model.DeletedUsersFromChat{DeletedUsers: deleted}, nil
}

// --- Private Helpers ---

func (uc *ChatUsecase) decorateDialog(ctx context.Context, chat *model.Chat, me uuid.UUID) {
	users, err := uc.chatRepo.GetUsersFromChat(ctx, chat.ID)
	if err != nil {
		zap.L().Warn("decorateDialog: GetUsersFromChat failed", zap.Error(err))
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
		if model.ChatType(c.Type) == model.ChatTypeDialog {
			users, err := uc.chatRepo.GetUsersFromChat(ctx, c.ID)
			if err != nil {
				continue
			}
			for _, u := range users {
				if u.Username == otherUsername {
					info, err := uc.GetChatInfo(ctx, me, c.ID)
					if err == nil {
						return info, true
					}
				}
			}
		}
	}
	return nil, false
}

func (uc *ChatUsecase) addDialogUsers(ctx context.Context, me uuid.UUID, otherUsername string, chatID uuid.UUID) {
	// owner is the one who started
	uc.chatRepo.AddUserToChatByID(ctx, me, string(model.RoleOwner), chatID)
	// other user is a member
	err := uc.chatRepo.AddUserToChatByUsername(ctx, otherUsername, string(model.RoleMember), chatID)
	if err != nil {
		zap.L().Warn("addDialogUsers: AddUserByUsername failed", zap.String("user", otherUsername), zap.Error(err))
	}
}

func (uc *ChatUsecase) addGroupOwner(ctx context.Context, me uuid.UUID, chatID uuid.UUID) {
	err := uc.chatRepo.AddUserToChatByID(ctx, me, string(model.RoleOwner), chatID)
	if err != nil {
		zap.L().Warn("addGroupOwner: AddUserToChatByID failed", zap.Error(err))
	}
}

func (uc *ChatUsecase) addChannelOwner(ctx context.Context, ownerID uuid.UUID, chatID uuid.UUID) {
	err := uc.chatRepo.AddUserToChatByID(ctx, ownerID, string(model.RoleOwner), chatID)
	if err != nil {
		zap.L().Warn("addChannelOwner: AddUserToChatByID failed", zap.String("chatID", chatID.String()), zap.Error(err))
	}
}

func (uc *ChatUsecase) extractUserIDs(users []model.UserInChat) []uuid.UUID {
	ids := make([]uuid.UUID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}
	return ids
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
