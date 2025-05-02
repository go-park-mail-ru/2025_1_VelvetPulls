package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/model"
	"github.com/google/uuid"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) SearchMessages(
	ctx context.Context,
	chatID uuid.UUID,
	query string,
	limit int,
	offset int,
) ([]model.Message, int, error) {
	// Поиск сообщений
	querySQL := `
        SELECT 
            m.id,
            m.body,
            m.user_id,
            m.sent_at,
            u.username
        FROM message m
        JOIN public.user u ON m.user_id = u.id
        WHERE m.chat_id = $1 
            AND m.body ILIKE $2
        ORDER BY sent_at DESC
        LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, querySQL, chatID, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, 0, ErrSearchMessages
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var m model.Message
		err := rows.Scan(
			&m.ID,
			&m.Body,
			&m.UserID,
			&m.SentAt,
			&m.Username,
		)
		if err != nil {
			return nil, 0, err
		}
		messages = append(messages, m)
	}

	// Получение общего количества
	var total int
	countSQL := `
		SELECT COUNT(*) 
		FROM message 
		WHERE chat_id = $1 
			AND to_tsvector('russian', body) 
			@@ plainto_tsquery('russian', $2)`
	err = r.db.QueryRowContext(ctx, countSQL, chatID, "%"+query+"%").Scan(&total)
	if err != nil {
		return nil, 0, ErrSearchMessagesCount
	}

	return messages, total, nil
}
