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
	AddUserToChatByID(ctx context.Context, userID uuid.UUID, userRole string, chatID uuid.UUID) error
	AddUserToChatByUsername(ctx context.Context, username string, userRole string, chatID uuid.UUID) error
	GetUserRoleInChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (string, error)
	GetUsersFromChat(ctx context.Context, chatId uuid.UUID) ([]model.UserInChat, error)
	RemoveUserFromChatByUsername(ctx context.Context, username string, chatID uuid.UUID) error
	RemoveUserFromChatByID(ctx context.Context, userID, chatID uuid.UUID) error
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) IChatRepo {
	return &chatRepository{db: db}
}

// --- Chat methods ---
func (r *chatRepository) GetChats(ctx context.Context, userID uuid.UUID) ([]model.Chat, uuid.UUID, error) {
	query := `
		SELECT 
			c.id, c.avatar_path, c.type, c.title,
			m.id, m.user_id, m.body, m.sent_at,
			(
				SELECT COUNT(*) 
				FROM user_chat uc2 
				WHERE uc2.chat_id = c.id
			) AS count_users
		FROM chat c
		JOIN user_chat uc ON c.id = uc.chat_id
		LEFT JOIN LATERAL (
			SELECT m.id, m.user_id, m.body, m.sent_at
			FROM message m
			WHERE m.chat_id = c.id
			ORDER BY m.sent_at DESC
			LIMIT 1
		) m ON true
		WHERE uc.user_id = $1
		ORDER BY m.sent_at DESC NULLS LAST
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, uuid.Nil, err
	}
	defer rows.Close()

	var chats []model.Chat
	var lastChatID uuid.UUID

	for rows.Next() {
		var chat model.Chat
		var msgID sql.NullString
		var msgUserID sql.NullString
		var msgBody sql.NullString
		var msgSentAt sql.NullTime

		err := rows.Scan(
			&chat.ID, &chat.AvatarPath, &chat.Type, &chat.Title,
			&msgID, &msgUserID, &msgBody, &msgSentAt, &chat.CountUsers,
		)
		if err != nil {
			return nil, uuid.Nil, err
		}

		if msgID.Valid && msgUserID.Valid && msgBody.Valid && msgSentAt.Valid {
			msgUUID, err1 := uuid.Parse(msgID.String)
			userUUID, err2 := uuid.Parse(msgUserID.String)
			if err1 == nil && err2 == nil {
				chat.LastMessage = &model.LastMessage{
					ID:     msgUUID,
					UserID: userUUID,
					Body:   msgBody.String,
					SentAt: msgSentAt.Time,
				}
			}
		}

		chats = append(chats, chat)
		lastChatID = chat.ID
	}

	if err := rows.Err(); err != nil {
		return nil, uuid.Nil, err
	}

	return chats, lastChatID, nil
}

func (r *chatRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	query := `SELECT id, avatar_path, type, title FROM chat WHERE id = $1`
	var chat model.Chat
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(&chat.ID, &chat.AvatarPath, &chat.Type, &chat.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChatNotFound
		}
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) CreateChat(ctx context.Context, create *model.CreateChat) (uuid.UUID, string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	var (
		avatarPath string
		args       []interface{}
	)

	if create.Avatar != nil {
		avatarPath = "./uploads/chats/" + uuid.New().String() + ".png"
		args = append(args, avatarPath)
	} else {
		args = append(args, nil)
	}
	args = append(args, create.Type, create.Title)

	query := `
		INSERT INTO chat (avatar_path, type, title, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	var chatID uuid.UUID
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&chatID); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return uuid.Nil, "", ErrRecordAlreadyExists
		}
		logger.Error("CreateChat failed", zap.Error(err))
		fmt.Println(err)
		return uuid.Nil, "", ErrDatabaseOperation
	}

	return chatID, avatarPath, nil
}

func (r *chatRepository) UpdateChat(ctx context.Context, update *model.UpdateChat) (string, string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	if update == nil || update.ID == uuid.Nil {
		return "", "", ErrInvalidUUID
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("BeginTx failed", zap.Error(err))
		return "", "", ErrDatabaseOperation
	}

	var (
		updates       []string
		args          []interface{}
		avatarOldPath string
		avatarNewPath string
		argIndex      = 1
	)

	if update.Avatar != nil {
		query := `SELECT avatar_path FROM chat WHERE id = $1 FOR UPDATE`
		var oldAvatar *string
		if err := tx.QueryRowContext(ctx, query, update.ID).Scan(&oldAvatar); err != nil && !errors.Is(err, sql.ErrNoRows) {
			rollbackTx(logger, tx)
			return "", "", ErrDatabaseOperation
		}
		if oldAvatar != nil {
			avatarOldPath = *oldAvatar
		}
		avatarNewPath = "./uploads/chats/" + uuid.New().String() + ".png"
		updates = append(updates, fmt.Sprintf("avatar_path = $%d", argIndex))
		args = append(args, avatarNewPath)
		argIndex++
	}

	if update.Title != nil {
		if *update.Title == "" {
			rollbackTx(logger, tx)
			return "", "", ErrEmptyField
		}
		updates = append(updates, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *update.Title)
		argIndex++
	}

	if len(updates) == 0 {
		rollbackTx(logger, tx)
		return "", "", nil
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	query := fmt.Sprintf("UPDATE chat SET %s WHERE id = $%d", strings.Join(updates, ", "), argIndex)
	args = append(args, update.ID)

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		rollbackTx(logger, tx)
		if strings.Contains(err.Error(), "duplicate key value") {
			return "", "", ErrRecordAlreadyExists
		}
		return "", "", ErrDatabaseOperation
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		rollbackTx(logger, tx)
		return "", "", ErrUpdateFailed
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Commit failed", zap.Error(err))
		return "", "", ErrDatabaseOperation
	}

	return avatarNewPath, avatarOldPath, nil
}

func (r *chatRepository) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chat WHERE id = $1`, chatID)
	return err
}

func (r *chatRepository) AddUserToChatByID(ctx context.Context, userID uuid.UUID, userRole string, chatID uuid.UUID) error {
	query := `INSERT INTO user_chat (user_id, chat_id, user_role, joined_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.db.ExecContext(ctx, query, userID, chatID, userRole)
	return err
}

func (r *chatRepository) AddUserToChatByUsername(ctx context.Context, username, userRole string, chatID uuid.UUID) error {
	var userID uuid.UUID
	if err := r.db.QueryRowContext(ctx, `SELECT id FROM public.user WHERE username = $1`, username).Scan(&userID); err != nil {
		return err
	}
	return r.AddUserToChatByID(ctx, userID, userRole, chatID)
}

func (r *chatRepository) GetUserRoleInChat(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) (string, error) {
	var role string
	err := r.db.QueryRowContext(ctx, `SELECT user_role FROM user_chat WHERE user_id = $1 AND chat_id = $2`, userID, chatID).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return role, err
}

func (r *chatRepository) GetUsersFromChat(ctx context.Context, chatID uuid.UUID) ([]model.UserInChat, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.id, u.username, u.name, u.avatar_path, uc.user_role
		FROM public.user u
		JOIN user_chat uc ON u.id = uc.user_id
		WHERE uc.chat_id = $1`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserInChat
	for rows.Next() {
		var user model.UserInChat
		if err := rows.Scan(&user.ID, &user.Username, &user.Name, &user.AvatarPath, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *chatRepository) RemoveUserFromChatByUsername(ctx context.Context, username string, chatID uuid.UUID) error {
	var userID uuid.UUID
	if err := r.db.QueryRowContext(ctx, `SELECT id FROM public.user WHERE username = $1`, username).Scan(&userID); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_chat WHERE user_id = $1 AND chat_id = $2`, userID, chatID)
	return err
}

func (r *chatRepository) RemoveUserFromChatByID(ctx context.Context, userID, chatID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_chat WHERE user_id = $1 AND chat_id = $2`, userID, chatID)
	return err
}
