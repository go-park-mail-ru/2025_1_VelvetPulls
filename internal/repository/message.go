package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type IMessageRepo interface {
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]model.Message, error)
	CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error)
}

type messageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) IMessageRepo {
	return &messageRepo{db: db}
}

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
	ORDER BY m.sent_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("GetMessages: query failed: %w", err)
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
			return nil, fmt.Errorf("GetMessages: scan failed: %w", err)
		}
		msg.Body = utils.SanitizeRichText(msg.Body)
		msg.Username = utils.SanitizeString(msg.Username)

		if parentMsgID.Valid {
			id, err := uuid.Parse(parentMsgID.String)
			if err == nil {
				msg.ParentMessageID = &id
			}
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetMessages: rows iteration failed: %w", err)
	}

	return messages, nil
}

func (r *messageRepo) getMessage(ctx context.Context, id uuid.UUID) (*model.Message, error) {
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
	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("getMessage: failed to get message: %w", err)
	}

	msg.Body = utils.SanitizeRichText(msg.Body)
	msg.Username = utils.SanitizeString(msg.Username)

	if parentMsgID.Valid {
		id, err := uuid.Parse(parentMsgID.String)
		if err == nil {
			msg.ParentMessageID = &id
		}
	}

	return &msg, nil
}

func (r *messageRepo) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	cleanBody := utils.SanitizeRichText(message.Body)
	if cleanBody == "" {
		return nil, ErrEmptyMessage
	}
	query := `
	INSERT INTO message (user_id, chat_id, body)
	VALUES ($1, $2, $3) 
	RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		message.UserID,
		message.ChatID,
		cleanBody,
	).Scan(
		&message.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateMessage: insert failed: %w", err)
	}
	fmt.Println(message.ID)
	messageOut, err := r.getMessage(ctx, message.ID)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return messageOut, nil
}
