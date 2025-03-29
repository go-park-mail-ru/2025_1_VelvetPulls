package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type IChatRepo interface {
	GetChatsByUserID(ctx context.Context) ([]model.Chat, error)
}

type chatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) IChatRepo {
	return &chatRepo{db: db}
}

func (r *chatRepo) GetChatsByUserID(ctx context.Context) ([]model.Chat, error) {
	userID, ok := ctx.Value(utils.USER_ID_KEY).(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("user_id not found in context")
	}

	query := `SELECT 
		c.id, 
		c.avatar_path, 
		c.type, 
		c.title,
		c.created_at, 
		c.updated_at 
	FROM public.chat c
	JOIN public.user_chat uc ON c.id = uc.chat_id
	WHERE uc.user_id = $1`
	ID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID format: %w", err)
	}

	rows, err := r.db.Query(query, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.Chat
	for rows.Next() {
		var chat model.Chat
		if err := rows.Scan(
			&chat.ID,
			&chat.AvatarPath,
			&chat.Type,
			&chat.Title,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
