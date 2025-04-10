package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IChatUsecase interface {
	GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, error)
	GetChatInfo(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (*model.ChatInfo, error)
	CreateChat(ctx context.Context, userID uuid.UUID, chat *model.CreateChat) (*model.ChatInfo, error)
	UpdateChat(ctx context.Context, userID uuid.UUID, chat *model.UpdateChat) (*model.ChatInfo, error)
	DeleteChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error
	AddUsersIntoChat(ctx context.Context, userID uuid.UUID, userIDs []uuid.UUID, chatID uuid.UUID) (*model.AddedUsersIntoChat, error)
	DeleteUserFromChat(ctx context.Context, userID uuid.UUID, userIDsDelete []uuid.UUID, chatID uuid.UUID) (*model.DeletedUsersFromChat, error)
}

type ChatUsecase struct {
	chatRepo repository.IChatRepo
}

func NewChatUsecase(chatRepo repository.IChatRepo) IChatUsecase {
	return &ChatUsecase{chatRepo: chatRepo}
}

func (uc *ChatUsecase) GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Getting chats for user", zap.String("userID", userID.String()))

	chats, _, err := uc.chatRepo.GetChats(ctx, userID)
	if err != nil {
		logger.Error("Failed to get chats from repository", zap.Error(err))
		return nil, err
	}

	for i := range chats {
		if chats[i].Type == "dialog" {
			users, err := uc.chatRepo.GetUsersFromChat(ctx, chats[i].ID)
			if err != nil {
				logger.Error("Failed to get users from chat",
					zap.String("chatID", chats[i].ID.String()),
					zap.Error(err))
				return nil, err
			}

			if len(users) > 0 {
				chats[i].Title = users[0].Username
				chats[i].AvatarPath = users[0].AvatarPath
			}
		}
	}

	logger.Info("Successfully retrieved chats", zap.Int("count", len(chats)))
	return chats, nil
}

func (uc *ChatUsecase) GetChatInfo(ctx context.Context, userID, chatID uuid.UUID) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Getting chat info",
		zap.String("userID", userID.String()),
		zap.String("chatID", chatID.String()))

	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("Failed to get user role in chat", zap.Error(err))
		return nil, err
	}
	if role != "owner" && role != "member" {
		logger.Warn("Permission denied for user in chat",
			zap.String("role", role))
		return nil, ErrPermissionDenied
	}

	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		logger.Error("Failed to get chat by ID", zap.Error(err))
		return nil, err
	}

	users, err := uc.chatRepo.GetUsersFromChat(ctx, chatID)
	if err != nil {
		logger.Error("Failed to get users from chat", zap.Error(err))
		return nil, err
	}

	if chat.Type == "dialog" {
		if len(users) > 0 {
			chat.Title = users[0].Username
			chat.AvatarPath = users[0].AvatarPath
		}
	}

	logger.Info("Successfully retrieved chat info")
	return &model.ChatInfo{
		ID:         chat.ID,
		AvatarPath: chat.AvatarPath,
		Type:       chat.Type,
		Title:      chat.Title,
		CountUsers: len(users),
	}, nil
}

func (uc *ChatUsecase) CreateChat(ctx context.Context, userID uuid.UUID, chat *model.CreateChat) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Creating chat",
		zap.String("userID", userID.String()),
		zap.String("chatType", chat.Type))

	if err := chat.Validate(); err != nil {
		logger.Warn("Chat validation failed", zap.Error(err))
		return nil, err
	}

	if chat.Avatar != nil {
		if !utils.IsImageFile(*chat.Avatar) {
			logger.Warn("Invalid avatar file type")
			return nil, utils.ErrNotImage
		}
	}

	if chat.Type == "dialog" {
		logger.Info("Checking for existing dialog")
		userChats, _, err := uc.chatRepo.GetChats(ctx, userID)
		if err != nil {
			logger.Error("Failed to get user chats", zap.Error(err))
			return nil, err
		}

		for _, userChat := range userChats {
			if userChat.Type == "dialog" {
				users, err := uc.chatRepo.GetUsersFromChat(ctx, userChat.ID)
				if err != nil {
					logger.Warn("Failed to get users from chat, skipping",
						zap.String("chatID", userChat.ID.String()),
						zap.Error(err))
					continue
				}

				for _, u := range users {
					if u.ID == chat.DialogUser {
						logger.Info("Existing dialog found, returning it")
						return uc.GetChatInfo(ctx, userID, userChat.ID)
					}
				}
			}
		}
	}

	chatID, avatarNewURL, err := uc.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		logger.Error("Failed to create chat in repository", zap.Error(err))
		return nil, ErrChatCreationFailed
	}

	logger.Info("Chat created, adding users", zap.String("chatID", chatID.String()))
	switch chat.Type {
	case "dialog":
		if err := uc.chatRepo.AddUserToChat(ctx, userID, "owner", chatID); err != nil {
			logger.Error("Failed to add owner to dialog", zap.Error(err))
			return nil, ErrAddOwnerToDialog
		}

		if err := uc.chatRepo.AddUserToChat(ctx, chat.DialogUser, "owner", chatID); err != nil {
			logger.Error("Failed to add participant to dialog", zap.Error(err))
			return nil, ErrAddParticipantToDialog
		}

	case "group":
		if chat.Avatar != nil {
			if !utils.IsImageFile(*chat.Avatar) {
				logger.Warn("Invalid avatar file type for group")
				return nil, utils.ErrNotImage
			}
		}

		if avatarNewURL != "" && chat.Avatar != nil {
			if err := utils.RewritePhoto(*chat.Avatar, avatarNewURL); err != nil {
				logger.Error("Failed to rewrite avatar photo", zap.Error(err))
				return nil, err
			}
		}

		if err := uc.chatRepo.AddUserToChat(ctx, userID, "owner", chatID); err != nil {
			logger.Error("Failed to add owner to group", zap.Error(err))
			return nil, ErrAddOwnerToGroup
		}
	}

	logger.Info("Chat successfully created")
	return uc.GetChatInfo(ctx, userID, chatID)
}

// ... (остальные методы аналогично с добавлением логгера)

func (uc *ChatUsecase) UpdateChat(ctx context.Context, userID uuid.UUID, chat *model.UpdateChat) (*model.ChatInfo, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Updating chat",
		zap.String("userID", userID.String()),
		zap.String("chatID", chat.ID.String()))

	chatFromDB, err := uc.chatRepo.GetChatByID(ctx, chat.ID)
	if err != nil {
		logger.Error("Failed to get chat from DB", zap.Error(err))
		return nil, err
	}

	if chatFromDB.Type == "dialog" {
		logger.Warn("Attempt to update dialog")
		return nil, ErrDialogUpdateForbidden
	}

	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chat.ID)
	if err != nil {
		logger.Error("Failed to get user role in chat", zap.Error(err))
		return nil, err
	}
	if role != "owner" {
		logger.Warn("User is not owner, can't modify chat",
			zap.String("role", role))
		return nil, ErrOnlyOwnerCanModify
	}

	if chat.Avatar != nil {
		if !utils.IsImageFile(*chat.Avatar) {
			logger.Warn("Invalid avatar file type")
			return nil, utils.ErrNotImage
		}
	}

	avatarNewURL, avatarOldURL, err := uc.chatRepo.UpdateChat(ctx, chat)
	if err != nil {
		logger.Error("Failed to update chat in repository", zap.Error(err))
		return nil, err
	}

	if avatarNewURL != "" && chat.Avatar != nil {
		if err := utils.RewritePhoto(*chat.Avatar, avatarNewURL); err != nil {
			logger.Error("Failed to rewrite avatar photo", zap.Error(err))
			return nil, err
		}
		if avatarOldURL != "" {
			go func() {
				if err := utils.RemovePhoto(avatarOldURL); err != nil {
					logger.Warn("Failed to remove old avatar photo",
						zap.String("path", avatarOldURL),
						zap.Error(err))
				}
			}()
		}
	}

	logger.Info("Chat successfully updated")
	return uc.GetChatInfo(ctx, userID, chat.ID)
}

func (uc *ChatUsecase) DeleteChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Deleting chat",
		zap.String("userID", userID.String()),
		zap.String("chatID", chatID.String()))

	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("Failed to get user role in chat", zap.Error(err))
		return err
	}
	if role != "owner" {
		logger.Warn("User is not owner, can't delete chat",
			zap.String("role", role))
		return ErrOnlyOwnerCanDelete
	}

	if err := uc.chatRepo.DeleteChat(ctx, chatID); err != nil {
		logger.Error("Failed to delete chat from repository", zap.Error(err))
		return err
	}

	logger.Info("Chat successfully deleted")
	return nil
}

func (uc *ChatUsecase) AddUsersIntoChat(ctx context.Context, userID uuid.UUID, userIDs []uuid.UUID, chatID uuid.UUID) (*model.AddedUsersIntoChat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Adding users to chat",
		zap.String("userID", userID.String()),
		zap.String("chatID", chatID.String()),
		zap.Int("usersCount", len(userIDs)))

	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		logger.Error("Failed to get chat by ID", zap.Error(err))
		return nil, err
	}

	if chat.Type == "dialog" {
		logger.Warn("Attempt to add users to dialog")
		return nil, ErrDialogAddUsers
	}

	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("Failed to get user role in chat", zap.Error(err))
		return nil, err
	}
	if role != "owner" {
		logger.Warn("User is not owner, can't add users",
			zap.String("role", role))
		return nil, ErrOnlyOwnerCanAddUsers
	}

	var added, notAdded []uuid.UUID
	for _, uid := range userIDs {
		if err := uc.chatRepo.AddUserToChat(ctx, uid, "member", chatID); err != nil {
			logger.Warn("Failed to add user to chat",
				zap.String("targetUserID", uid.String()),
				zap.Error(err))
			notAdded = append(notAdded, uid)
		} else {
			logger.Info("User added to chat",
				zap.String("targetUserID", uid.String()))
			added = append(added, uid)
		}
	}

	logger.Info("Users addition completed",
		zap.Int("added", len(added)),
		zap.Int("notAdded", len(notAdded)))
	return &model.AddedUsersIntoChat{AddedUsers: added, NotAddedUsers: notAdded}, nil
}

func (uc *ChatUsecase) DeleteUserFromChat(ctx context.Context, userID uuid.UUID, userIDsDelete []uuid.UUID, chatID uuid.UUID) (*model.DeletedUsersFromChat, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Deleting users from chat",
		zap.String("userID", userID.String()),
		zap.String("chatID", chatID.String()),
		zap.Int("usersCount", len(userIDsDelete)))

	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		logger.Error("Failed to get chat by ID", zap.Error(err))
		return nil, err
	}

	if chat.Type == "dialog" {
		logger.Warn("Attempt to delete users from dialog")
		return nil, ErrDialogDeleteUsers
	}

	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("Failed to get user role in chat", zap.Error(err))
		return nil, err
	}
	if role != "owner" {
		logger.Warn("User is not owner, can't delete users",
			zap.String("role", role))
		return nil, ErrOnlyOwnerCanDeleteUsers
	}

	var deleted []uuid.UUID
	for _, uid := range userIDsDelete {
		if err := uc.chatRepo.RemoveUserFormChat(ctx, uid, chatID); err == nil {
			logger.Info("User removed from chat",
				zap.String("targetUserID", uid.String()))
			deleted = append(deleted, uid)
		} else {
			logger.Warn("Failed to remove user from chat",
				zap.String("targetUserID", uid.String()),
				zap.Error(err))
		}
	}

	logger.Info("Users deletion completed",
		zap.Int("deleted", len(deleted)))
	return &model.DeletedUsersFromChat{DeletedUsers: deleted}, nil
}
