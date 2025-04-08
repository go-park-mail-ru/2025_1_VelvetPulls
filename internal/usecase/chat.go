package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
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
	chats, _, err := uc.chatRepo.GetChats(ctx, userID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (uc *ChatUsecase) GetChatInfo(ctx context.Context, userID, chatID uuid.UUID) (*model.ChatInfo, error) {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}

	if role == "" {
		return nil, fmt.Errorf("permission denied")
	}

	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	users, err := uc.chatRepo.GetUsersFromChat(ctx, chatID)
	if err != nil {
		return nil, err
	}

	return &model.ChatInfo{
		ID:         chat.ID,
		AvatarPath: chat.AvatarPath,
		Type:       chat.Type,
		Title:      chat.Title,
		CountUsers: len(users),
	}, nil
}

func (uc *ChatUsecase) CreateChat(ctx context.Context, userID uuid.UUID, chat *model.CreateChat) (*model.ChatInfo, error) {
	// Валидация типа чата
	if err := chat.Validate(); err != nil {
		return nil, err
	}
	// Проверка аватара
	if chat.Avatar != nil {
		if !utils.IsImageFile(*chat.Avatar) {
			return nil, utils.ErrNotImage
		}
	}

	// Создаем чат в репозитории
	chatID, _, err := uc.chatRepo.CreateChat(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	// Добавляем создателя чата
	if err := uc.chatRepo.AddUserToChat(ctx, userID, "owner", chatID); err != nil {
		return nil, fmt.Errorf("failed to add owner to chat: %w", err)
	}

	return uc.GetChatInfo(ctx, userID, chatID)
}

func (uc *ChatUsecase) UpdateChat(ctx context.Context, userID uuid.UUID, chat *model.UpdateChat) (*model.ChatInfo, error) {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chat.ID)
	if err != nil {
		return nil, err
	}

	if role != "owner" {
		return nil, fmt.Errorf("only chat owner can delete users")
	}

	if chat.Avatar != nil {
		if !utils.IsImageFile(*chat.Avatar) {
			return nil, utils.ErrNotImage
		}
	}

	avatarNewURL, avatarOldURL, err := uc.chatRepo.UpdateChat(ctx, chat)
	if err != nil {
		return nil, err
	}

	// Если есть новый аватар, сохраняем его и удаляем старый
	if avatarNewURL != "" && chat.Avatar != nil {
		if err := utils.RewritePhoto(*chat.Avatar, avatarNewURL); err != nil {
			// logger.Error("Error rewriting photo")
			return nil, err
		}
		if avatarOldURL != "" {
			go func() {
				if err := utils.RemovePhoto(avatarOldURL); err != nil {
					// logger.Error("Error removing old avatar")
				}
			}()
		}
	}

	return uc.GetChatInfo(ctx, chat.ID, chat.ID)
}

func (uc *ChatUsecase) DeleteChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if role != "owner" {
		return fmt.Errorf("only chat owner can delete chat")
	}

	return uc.chatRepo.DeleteChat(ctx, chatID)
}

func (uc *ChatUsecase) AddUsersIntoChat(ctx context.Context, userID uuid.UUID, userIDs []uuid.UUID, chatID uuid.UUID) (*model.AddedUsersIntoChat, error) {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}
	if role != "owner" {
		return nil, fmt.Errorf("only chat owner can add users")
	}

	var added, notAdded []uuid.UUID
	for _, uid := range userIDs {
		if err := uc.chatRepo.AddUserToChat(ctx, uid, "member", chatID); err != nil {
			notAdded = append(notAdded, uid)
		} else {
			added = append(added, uid)
		}
	}
	return &model.AddedUsersIntoChat{AddedUsers: added, NotAddedUsers: notAdded}, nil
}

func (uc *ChatUsecase) DeleteUserFromChat(ctx context.Context, userID uuid.UUID, userIDsDelete []uuid.UUID, chatID uuid.UUID) (*model.DeletedUsersFromChat, error) {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}
	if role != "owner" {
		return nil, fmt.Errorf("only chat owner can delete users")
	}

	var deleted []uuid.UUID
	for _, uid := range userIDsDelete {
		if err := uc.chatRepo.RemoveUserFormChat(ctx, uid, chatID); err == nil {
			deleted = append(deleted, uid)
		}
	}
	return &model.DeletedUsersFromChat{DeletedUsers: deleted}, nil
}
