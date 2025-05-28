package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/google/uuid"
)

type ChatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) SearchUserChats(ctx context.Context, userID uuid.UUID, query string) ([]model.Chat, error) {
	baseQuery := `
        SELECT 
            c.id, 
            c.type, 
            c.title, 
            c.avatar_path,
            c.created_at,
            c.updated_at,
			uc.send_notifications,
            m.id AS last_message_id,
            m.user_id AS last_message_user_id,
            m.body AS last_message_body,
            m.sent_at AS last_message_sent_at
        FROM chat c
        LEFT JOIN (
            SELECT 
                *,
                ROW_NUMBER() OVER (PARTITION BY chat_id ORDER BY sent_at DESC) as rn
            FROM message
        ) m ON c.id = m.chat_id AND m.rn = 1
        JOIN user_chat uc ON c.id = uc.chat_id
        WHERE uc.user_id = $1
    `

	args := []interface{}{userID}
	paramCount := 2

	if query != "" {
		baseQuery += fmt.Sprintf(` AND c.title ILIKE $%d`, paramCount)
		args = append(args, "%"+query+"%")
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, ErrSearchUserChats
	}
	defer rows.Close()

	var chats []model.Chat
	for rows.Next() {
		var c model.Chat
		var sendNotifications bool
		var lastMessageID sql.NullString
		var lastMessageUserID sql.NullString
		var lastMessageBody sql.NullString
		var lastMessageSentAt sql.NullTime

		err := rows.Scan(
			&c.ID,
			&c.Type,
			&c.Title,
			&c.AvatarPath,
			&c.CreatedAt,
			&c.UpdatedAt,
			&sendNotifications,
			&lastMessageID,
			&lastMessageUserID,
			&lastMessageBody,
			&lastMessageSentAt,
		)

		if err != nil {
			return nil, err
		}

		if lastMessageID.Valid {
			c.LastMessage = &model.LastMessage{
				ID:     uuid.MustParse(lastMessageID.String),
				UserID: uuid.MustParse(lastMessageUserID.String),
				Body:   lastMessageBody.String,
				SentAt: lastMessageSentAt.Time,
			}
		}

		chats = append(chats, c)
	}

	return chats, nil
}

func (r *ChatRepo) SearchGlobalChannels(ctx context.Context, query string) ([]model.Chat, error) {
	baseQuery := `
	SELECT DISTINCT ON (c.id)
		c.id, 
		c.type, 
		c.title, 
		c.avatar_path,
		c.created_at,
		c.updated_at,
		uc.send_notifications,
		m.id AS last_message_id,
		m.user_id AS last_message_user_id,
		m.body AS last_message_body,
		m.sent_at AS last_message_sent_at
	FROM chat c
	LEFT JOIN (
		SELECT 
			*,
			ROW_NUMBER() OVER (PARTITION BY chat_id ORDER BY sent_at DESC) as rn
		FROM message
	) m ON c.id = m.chat_id AND m.rn = 1
	JOIN user_chat uc ON c.id = uc.chat_id
	WHERE c.type = 'channel'
	`

	args := []interface{}{}
	paramCount := 1

	if query != "" {
		baseQuery += fmt.Sprintf(` AND c.title ILIKE $%d`, paramCount)
		args = append(args, "%"+query+"%")
	}

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, ErrSearchUserChats
	}
	defer rows.Close()

	var chats []model.Chat
	for rows.Next() {
		var c model.Chat
		var sendNotifications bool
		var lastMessageID sql.NullString
		var lastMessageUserID sql.NullString
		var lastMessageBody sql.NullString
		var lastMessageSentAt sql.NullTime

		err := rows.Scan(
			&c.ID,
			&c.Type,
			&c.Title,
			&c.AvatarPath,
			&c.CreatedAt,
			&c.UpdatedAt,
			&sendNotifications,
			&lastMessageID,
			&lastMessageUserID,
			&lastMessageBody,
			&lastMessageSentAt,
		)

		if err != nil {
			return nil, err
		}

		if lastMessageID.Valid {
			c.LastMessage = &model.LastMessage{
				ID:     uuid.MustParse(lastMessageID.String),
				UserID: uuid.MustParse(lastMessageUserID.String),
				Body:   lastMessageBody.String,
				SentAt: lastMessageSentAt.Time,
			}
		}

		chats = append(chats, c)
	}

	return chats, nil
}

func (r *ChatRepo) GetUsersFromChat(ctx context.Context, chatID uuid.UUID) ([]model.UserInChat, error) {
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
