package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

// TODO: Переделать под новую структуру бд
type IChatRepo interface {
	GetChatsByUserID(ctx context.Context, username string) ([]model.Chat, error)
}

type chatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) IChatRepo {
	return &chatRepo{db: db}
}

// TODO: поправить под user_id
func (r *chatRepo) GetChatsByUserID(ctx context.Context, userID string) ([]model.Chat, error) {
	rows, err := r.db.Query("SELECT id, owner_username, type, title, description, created_at, updated_at FROM chats WHERE owner_username = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.Chat
	for rows.Next() {
		var chat model.Chat
		if err := rows.Scan(&chat.ID, &chat.OwnerUsername, &chat.Type, &chat.Title, &chat.Description, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
