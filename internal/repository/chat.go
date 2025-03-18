package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

type IChatRepo interface {
	GetChatsByUsername(username string) ([]model.Chat, error)
	AddChat(chat *model.Chat) error
}

type chatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) IChatRepo {
	return &chatRepo{db: db}
}

func (r *chatRepo) GetChatsByUsername(username string) ([]model.Chat, error) {
	rows, err := r.db.Query("SELECT id, owner_username, type, title, description, created_at, updated_at FROM chats WHERE owner_username = $1", username)
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

func (r *chatRepo) AddChat(chat *model.Chat) error {
	query := "INSERT INTO chats (owner_username, type, title, description, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id"
	return r.db.QueryRow(query, chat.OwnerUsername, chat.Type, chat.Title, chat.Description).Scan(&chat.ID)
}
