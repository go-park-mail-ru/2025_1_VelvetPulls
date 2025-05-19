package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
)

type IMessageRepo interface {
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]model.Message, error)
	GetMessagesBefore(ctx context.Context, chatID uuid.UUID, beforeMessageID uuid.UUID) ([]model.Message, error)
	GetMessagesAfter(ctx context.Context, chatID uuid.UUID, beforeMessageID uuid.UUID) ([]model.Message, error)
	GetMessage(ctx context.Context, id uuid.UUID) (*model.Message, error)
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	UpdateMessage(ctx context.Context, messageID uuid.UUID, newBody string) (*model.Message, error)
	DeleteMessage(ctx context.Context, messageID uuid.UUID) (*model.Message, error)
}

type messageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) IMessageRepo {
	return &messageRepo{db: db}
}

var limit = 25

func (r *messageRepo) GetMessages(ctx context.Context, chatID uuid.UUID) ([]model.Message, error) {
	query := `
		SELECT 
			m.id,
			m.parent_message_id,
			m.chat_id,
			m.user_id,
			m.body,
			m.sent_at,
			m.is_redacted,
			u.username,
			u.avatar_path
		FROM message m
		JOIN public.user u ON m.user_id = u.id
		WHERE m.chat_id = $1
		ORDER BY m.sent_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, chatID, limit)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString

		if err := rows.Scan(
			&msg.ID,
			&parentMsgID,
			&msg.ChatID,
			&msg.UserID,
			&msg.Body,
			&msg.SentAt,
			&msg.IsRedacted,
			&msg.Username,
			&msg.AvatarPath,
		); err != nil {
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			id, err := uuid.Parse(parentMsgID.String)
			if err == nil {
				msg.ParentMessageID = &id
			}
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrDatabaseOperation
	}

	return messages, nil
}

func (r *messageRepo) GetMessagesBefore(ctx context.Context, chatID uuid.UUID, beforeMessageID uuid.UUID) ([]model.Message, error) {
	query := `
    WITH ref_message AS (
        SELECT sent_at, id FROM message WHERE id = $2
    )
    SELECT 
        m.id,
        m.parent_message_id,
        m.chat_id,
        m.user_id,
        m.body,
        m.sent_at,
        m.is_redacted,
        u.avatar_path,
        u.username
    FROM message m
    JOIN public.user u ON m.user_id = u.id
    JOIN ref_message r ON TRUE
    WHERE m.chat_id = $1
      AND (
        m.sent_at < r.sent_at
        OR (m.sent_at = r.sent_at AND m.id::text < r.id::text)
      )
    ORDER BY m.sent_at DESC, m.id DESC
    LIMIT $3
    `

	rows, err := r.db.QueryContext(ctx, query, chatID, beforeMessageID, limit)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString

		err := rows.Scan(
			&msg.ID,
			&parentMsgID,
			&msg.ChatID,
			&msg.UserID,
			&msg.Body,
			&msg.SentAt,
			&msg.IsRedacted,
			&avatarPath,
			&msg.Username,
		)
		if err != nil {
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			parentUUID, err := uuid.Parse(parentMsgID.String)
			if err == nil {
				msg.ParentMessageID = &parentUUID
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrDatabaseOperation
	}

	return messages, nil
}

func (r *messageRepo) GetMessagesAfter(ctx context.Context, chatID uuid.UUID, afterMessageID uuid.UUID) ([]model.Message, error) {
	query := `
    WITH ref_message AS (
        SELECT sent_at, id FROM message WHERE id = $2
    )
    SELECT 
        m.id,
        m.parent_message_id,
        m.chat_id,
        m.user_id,
        m.body,
        m.sent_at,
        m.is_redacted,
        u.avatar_path,
        u.username
    FROM message m
    JOIN public.user u ON m.user_id = u.id
    JOIN ref_message r ON TRUE
    WHERE m.chat_id = $1
      AND (
        m.sent_at > r.sent_at
        OR (m.sent_at = r.sent_at AND m.id::text > r.id::text)
      )
    ORDER BY m.sent_at ASC, m.id ASC
    LIMIT $3
    `

	rows, err := r.db.QueryContext(ctx, query, chatID, afterMessageID, limit)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString

		err := rows.Scan(
			&msg.ID,
			&parentMsgID,
			&msg.ChatID,
			&msg.UserID,
			&msg.Body,
			&msg.SentAt,
			&msg.IsRedacted,
			&avatarPath,
			&msg.Username,
		)
		if err != nil {
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			parentUUID, err := uuid.Parse(parentMsgID.String)
			if err == nil {
				msg.ParentMessageID = &parentUUID
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrDatabaseOperation
	}

	return messages, nil
}

func (r *messageRepo) GetMessage(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	query := `
	SELECT 
		m.id,
		m.parent_message_id,
		m.chat_id,
		m.user_id,
		m.body,
		m.sent_at,
		m.is_redacted,
		u.username,
		u.avatar_path
	FROM message m
	JOIN public.user u ON m.user_id = u.id
	WHERE m.id = $1
	`

	var msg model.Message
	var parentMsgID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&msg.ID,
		&parentMsgID,
		&msg.ChatID,
		&msg.UserID,
		&msg.Body,
		&msg.SentAt,
		&msg.IsRedacted,
		&msg.Username,
		&msg.AvatarPath,
	)
	if err != nil {
		return nil, ErrDatabaseOperation
	}

	if parentMsgID.Valid {
		id, err := uuid.Parse(parentMsgID.String)
		if err == nil {
			msg.ParentMessageID = &id
		}
	}

	return &msg, nil
}

func (r *messageRepo) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	query := `
	INSERT INTO message (user_id, chat_id, body)
	VALUES ($1, $2, $3) 
	RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		message.UserID,
		message.ChatID,
		message.Body,
	).Scan(&message.ID)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	messageOut, err := r.GetMessage(ctx, message.ID)
	if err != nil {
		return nil, err
	}
	return messageOut, nil
}

func (r *messageRepo) UpdateMessage(ctx context.Context, messageID uuid.UUID, newBody string) (*model.Message, error) {
	query := `
		UPDATE message
		SET body = $1, is_redacted = true
		WHERE id = $2
		RETURNING id
	`

	var updatedID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, newBody, messageID).Scan(&updatedID)
	if err != nil {
		return nil, ErrUpdateFailed
	}

	return r.GetMessage(ctx, updatedID)
}

func (r *messageRepo) DeleteMessage(ctx context.Context, messageID uuid.UUID) (*model.Message, error) {
	query := `
		DELETE FROM message
		WHERE id = $1
		RETURNING id
	`

	var deletedID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(&deletedID)
	if err != nil {
		return nil, ErrDatabaseOperation
	}

	return &model.Message{ID: deletedID}, nil
}
