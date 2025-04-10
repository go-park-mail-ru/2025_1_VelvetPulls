package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IChatRepo interface {
	GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, uuid.UUID, error)
	GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error)
	CreateChat(ctx context.Context, create *model.CreateChat) (uuid.UUID, string, error)
	UpdateChat(ctx context.Context, update *model.UpdateChat) (string, string, error)
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
	AddUserToChat(ctx context.Context, userID uuid.UUID, userRole string, chatID uuid.UUID) error
	GetUserRoleInChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (string, error)
	GetUsersFromChat(ctx context.Context, chatId uuid.UUID) ([]model.UserInChat, error)
	RemoveUserFormChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) IChatRepo {
	return &chatRepository{db: db}
}

func (r *chatRepository) GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, uuid.UUID, error) {
	var chats []model.Chat
	var lastChatID uuid.UUID
	query := `SELECT c.id, c.avatar_path, c.type, c.title, c.created_at, c.updated_at
			  FROM chat c
			  JOIN user_chat uc ON c.id = uc.chat_id
			  WHERE uc.user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, lastChatID, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.Chat
		if err := rows.Scan(&chat.ID, &chat.AvatarPath, &chat.Type, &chat.Title, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
			return nil, lastChatID, err
		}
		chats = append(chats, chat)
		lastChatID = chat.ID
	}
	return chats, lastChatID, nil
}

func (r *chatRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	var chat model.Chat
	query := `SELECT id, avatar_path, type, title, created_at, updated_at FROM chat WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, chatID)
	if err := row.Scan(&chat.ID, &chat.AvatarPath, &chat.Type, &chat.Title, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) CreateChat(ctx context.Context, create *model.CreateChat) (uuid.UUID, string, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	var avatarPath string
	var args []interface{}

	if create.Avatar != nil {
		logger.Info("Processing chat avatar")

		avatarDir := "./uploads/chats/"

		avatarPath = avatarDir + uuid.New().String() + ".png"
		args = append(args, avatarPath)
	} else {
		args = append(args, nil)
	}

	args = append(args, create.Type, create.Title)

	query := `
        INSERT INTO public.chat 
        (avatar_path, type, title, created_at, updated_at) 
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id
    `

	var chatID uuid.UUID
	logger.Info("Executing chat creation query")

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&chatID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("Chat with similar properties already exists")
			return uuid.UUID{}, "", ErrRecordAlreadyExists
		}
		logger.Error("Database operation failed", zap.Error(err))
		return uuid.UUID{}, "", ErrDatabaseOperation
	}

	logger.Info("Successfully created chat", zap.String("chatID", chatID.String()))
	return chatID, avatarPath, nil
}

func (r *chatRepository) UpdateChat(ctx context.Context, update *model.UpdateChat) (string, string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	if update == nil {
		logger.Warn("Update data is empty")
		return "", "", ErrEmptyField
	}
	if update.ID == uuid.Nil {
		logger.Warn("Invalid UUID provided")
		return "", "", ErrInvalidUUID
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction")
		return "", "", ErrDatabaseOperation
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var updates []string
	var args []interface{}
	argIndex := 1
	var avatarOldURL, avatarNewURL string

	if update.Avatar != nil {
		logger.Info("Updating chat avatar")

		var oldUrl *string
		err = tx.QueryRowContext(ctx,
			"SELECT avatar_path FROM public.chat WHERE id = $1 FOR UPDATE",
			update.ID,
		).Scan(&oldUrl)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("Failed to get current avatar path")
			return "", "", ErrDatabaseOperation
		}
		if oldUrl != nil {
			avatarOldURL = *oldUrl
		}

		avatarDir := "./uploads/chats/"
		avatarNewURL = avatarDir + uuid.New().String() + ".png"

		updates = append(updates, fmt.Sprintf("avatar_path = $%d", argIndex))
		args = append(args, avatarNewURL)
		argIndex++
	}

	if update.Title != nil {
		if *update.Title == "" {
			logger.Warn("Chat title cannot be empty")
			return "", "", ErrEmptyField
		}
		updates = append(updates, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *update.Title)
		argIndex++
	}

	if len(updates) == 0 {
		return "", "", nil
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	query := fmt.Sprintf(
		"UPDATE public.chat SET %s WHERE id = $%d",
		strings.Join(updates, ", "),
		argIndex,
	)
	args = append(args, update.ID)

	logger.Info("Executing chat update query inside transaction")
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("Attempt to update to existing values")
			return "", "", ErrRecordAlreadyExists
		}
		logger.Error("Database operation failed")
		return "", "", ErrDatabaseOperation
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("Failed to get affected rows count")
		return "", "", ErrDatabaseOperation
	}
	if rowsAffected == 0 {
		logger.Warn("No rows affected, update failed")
		return "", "", ErrUpdateFailed
	}

	return avatarNewURL, avatarOldURL, nil
}

func (r *chatRepository) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	query := `DELETE FROM chat WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, chatID)
	return err
}

func (r *chatRepository) AddUserToChat(ctx context.Context, userID uuid.UUID, userRole string, chatID uuid.UUID) error {
	query := `INSERT INTO user_chat (user_id, chat_id, user_role, joined_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.db.ExecContext(ctx, query, userID, chatID, userRole)
	return err
}

func (r *chatRepository) GetUserRoleInChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (string, error) {
	var role string
	query := `SELECT user_role FROM public.user_chat WHERE user_id = $1 AND chat_id = $2`
	row := r.db.QueryRowContext(ctx, query, userID, chatID)
	if err := row.Scan(&role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return role, nil
}

func (r *chatRepository) GetUsersFromChat(ctx context.Context, chatID uuid.UUID) ([]model.UserInChat, error) {
	var users []model.UserInChat
	query := `SELECT u.id, u.username, u.first_name, u.avatar_path, uc.user_role
			  FROM public.user u
			  JOIN public.user_chat uc ON u.id = uc.user_id
			  WHERE uc.chat_id = $1`
	rows, err := r.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.UserInChat
		if err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarPath, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *chatRepository) RemoveUserFormChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) error {
	query := `DELETE FROM user_chat WHERE user_id = $1 AND chat_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, chatID)
	return err
}

func (r *chatRepository) GetProfileAvatarAndName(ctx context.Context, userID uuid.UUID) (string, string, error) {
	query := `SELECT avatar_path, username FROM public.user WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)

	var userName string
	var avatarPath sql.NullString
	err := row.Scan(&avatarPath, &userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrUserNotFound
		}
		return "", "", ErrDatabaseOperation
	}

	actualAvatarPath := ""
	if avatarPath.Valid {
		actualAvatarPath = avatarPath.String
	}

	return actualAvatarPath, userName, nil
}
