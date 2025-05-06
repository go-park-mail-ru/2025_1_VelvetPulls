package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/google/uuid"
)

type ChatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) SearchUserChats(ctx context.Context, userID uuid.UUID, query string) ([]model.Chat, [][]model.UserInChat, error) {
	const querySQL = `
	WITH user_chats AS (
		SELECT 
			c.id,
			c.type,
			c.title,
			c.avatar_path,
			c.created_at,
			c.updated_at,
			uc.user_role,
			jsonb_agg(jsonb_build_object(
				'id', u.id::text,
				'username', u.username,
				'avatar_path', u.avatar_path
			)) AS participants
		FROM chat c
		JOIN user_chat uc ON c.id = uc.chat_id
		JOIN public.user u ON uc.user_id = u.id
		WHERE uc.user_id = $1
		GROUP BY c.id, uc.user_role
	)
	SELECT 
		uc.id,
		uc.type,
		uc.title,
		uc.avatar_path,
		uc.created_at,
		uc.updated_at,
		uc.user_role,
		uc.participants,
		m.id AS last_message_id,
		m.body AS last_message_body,
		m.sent_at AS last_message_sent_at,
		m.user_id AS last_message_user_id,
		u.username AS last_message_username
	FROM user_chats uc
	LEFT JOIN LATERAL (
		SELECT 
			m.id,
			m.body,
			m.sent_at,
			m.user_id
		FROM message m
		WHERE m.chat_id = uc.id 
		ORDER BY m.sent_at DESC 
		LIMIT 1
	) m ON true
	LEFT JOIN public.user u ON m.user_id = u.id
	WHERE (
		(uc.type IN ('group', 'channel') AND uc.title ILIKE $2)
		OR
		(uc.type = 'dialog' AND EXISTS (
			SELECT 1 
			FROM jsonb_array_elements(uc.participants) p 
			WHERE p->>'username' ILIKE $2
			AND p->>'id' != $1::text
		))
	)
	ORDER BY m.sent_at DESC NULLS LAST`

	rows, err := r.db.QueryContext(ctx, querySQL, userID, "%"+query+"%")
	if err != nil {
		return nil, nil, ErrSearchUserChats
	}
	defer rows.Close()

	var chats []model.Chat
	var participantsList [][]model.UserInChat

	for rows.Next() {
		var (
			c                 model.Chat
			lastMessageID     sql.NullString
			lastMessageBody   sql.NullString
			lastMessageSentAt sql.NullTime
			lastMessageUserID sql.NullString
			lastMessageUser   sql.NullString
			participantsJSON  []byte
		)

		err := rows.Scan(
			&c.ID,
			&c.Type,
			&c.Title,
			&c.AvatarPath,
			&c.CreatedAt,
			&c.UpdatedAt,
			&participantsJSON,
			&lastMessageID,
			&lastMessageBody,
			&lastMessageSentAt,
			&lastMessageUserID,
			&lastMessageUser,
		)

		if err != nil {
			return nil, nil, err
		}

		// Парсинг участников
		var participants []model.UserInChat
		if err := json.Unmarshal(participantsJSON, &participants); err != nil {
			return nil, nil, ErrUnmarshalParticipiants
		}

		// Обработка последнего сообщения
		if lastMessageID.Valid {
			c.LastMessage = &model.LastMessage{
				ID:       uuid.MustParse(lastMessageID.String),
				Body:     lastMessageBody.String,
				SentAt:   lastMessageSentAt.Time,
				UserID:   uuid.MustParse(lastMessageUserID.String),
				Username: lastMessageUser.String,
			}
		}

		chats = append(chats, c)
		participantsList = append(participantsList, participants)
	}

	return chats, participantsList, nil
}
