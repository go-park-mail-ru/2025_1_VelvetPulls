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
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM message WHERE chat_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, chatID).Scan(&totalCount); err != nil {
		return nil, ErrDatabaseOperation
	}

	if totalCount < limit {
		return nil, nil
	}

	queryBefore := `
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

	rows, err := r.db.QueryContext(ctx, queryBefore, chatID, beforeMessageID, limit)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message
	msgIDs := map[uuid.UUID]struct{}{}

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString

		if err := rows.Scan(
			&msg.ID,
			&parentMsgID,
			&msg.ChatID,
			&msg.UserID,
			&msg.Body,
			&msg.SentAt,
			&msg.IsRedacted,
			&avatarPath,
			&msg.Username,
		); err != nil {
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			if parentUUID, err := uuid.Parse(parentMsgID.String); err == nil {
				msg.ParentMessageID = &parentUUID
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		messages = append(messages, msg)
		msgIDs[msg.ID] = struct{}{}
	}

	if len(messages) >= limit {
		return messages, nil
	}

	// Нужно добрать
	queryFill := `
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
	WHERE m.chat_id = $1
	ORDER BY m.sent_at DESC, m.id DESC
	LIMIT $2
	`

	fillRows, err := r.db.QueryContext(ctx, queryFill, chatID, limit)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer fillRows.Close()

	for fillRows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString

		if err := fillRows.Scan(
			&msg.ID,
			&parentMsgID,
			&msg.ChatID,
			&msg.UserID,
			&msg.Body,
			&msg.SentAt,
			&msg.IsRedacted,
			&avatarPath,
			&msg.Username,
		); err != nil {
			return nil, ErrDatabaseScan
		}

		if _, exists := msgIDs[msg.ID]; exists {
			continue
		}

		if parentMsgID.Valid {
			if parentUUID, err := uuid.Parse(parentMsgID.String); err == nil {
				msg.ParentMessageID = &parentUUID
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		messages = append(messages, msg)

		if len(messages) >= limit {
			break
		}
	}

	return messages, nil
}

func (r *messageRepo) GetMessagesAfter(ctx context.Context, chatID uuid.UUID, afterMessageID uuid.UUID) ([]model.Message, error) {
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM message WHERE chat_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, chatID).Scan(&totalCount); err != nil {
		return nil, ErrDatabaseOperation
	}

	if totalCount < limit {
		return nil, nil
	}

	queryAfter := `
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

	rows, err := r.db.QueryContext(ctx, queryAfter, chatID, afterMessageID, limit)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message
	msgIDs := map[uuid.UUID]struct{}{}

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString

		if err := rows.Scan(
			&msg.ID,
			&parentMsgID,
			&msg.ChatID,
			&msg.UserID,
			&msg.Body,
			&msg.SentAt,
			&msg.IsRedacted,
			&avatarPath,
			&msg.Username,
		); err != nil {
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			if parentUUID, err := uuid.Parse(parentMsgID.String); err == nil {
				msg.ParentMessageID = &parentUUID
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		messages = append(messages, msg)
		msgIDs[msg.ID] = struct{}{}
	}

	if len(messages) >= limit {
		return messages, nil
	}

	// Добираем с начала (самые старые сообщения чата)
	queryFill := `
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
	WHERE m.chat_id = $1
	ORDER BY m.sent_at ASC, m.id ASC
	LIMIT $2
	`

	fillRows, err := r.db.QueryContext(ctx, queryFill, chatID, limit)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer fillRows.Close()

	for fillRows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString

		if err := fillRows.Scan(
			&msg.ID,
			&parentMsgID,
			&msg.ChatID,
			&msg.UserID,
			&msg.Body,
			&msg.SentAt,
			&msg.IsRedacted,
			&avatarPath,
			&msg.Username,
		); err != nil {
			return nil, ErrDatabaseScan
		}

		if _, exists := msgIDs[msg.ID]; exists {
			continue
		}

		if parentMsgID.Valid {
			if parentUUID, err := uuid.Parse(parentMsgID.String); err == nil {
				msg.ParentMessageID = &parentUUID
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		messages = append(messages, msg)

		if len(messages) >= limit {
			break
		}
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
