package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
)

const (
	limit                  = 25
	defaultMessageType     = "default"
	MessageWithPayloadType = "with_payload"
	stickerMessageType     = "sticker"
	filePayloadType        = "file"
	photoPayloadType       = "photo"
)

type IMessageRepo interface {
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]model.Message, error)
	GetMessagesBefore(ctx context.Context, chatID, beforeMessageID uuid.UUID) ([]model.Message, error)
	GetMessagesAfter(ctx context.Context, chatID, afterMessageID uuid.UUID) ([]model.Message, error)
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
			u.avatar_path,
			u.username,
			m.message_type,
			m.sticker_path
		FROM message m
		JOIN public.user u ON m.user_id = u.id
		WHERE m.chat_id = $1
		ORDER BY m.sent_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, chatID, limit)
	if err != nil {
		log.Println("get messages:", err)
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString
		var stickerPath sql.NullString

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
			&msg.MessageType,
			&stickerPath,
		)
		if err != nil {
			log.Println("scan message:", err)
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			id, err := uuid.Parse(parentMsgID.String)
			if err == nil {
				msg.ParentMessageID = &id
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		if stickerPath.Valid {
			msg.Sticker = stickerPath.String
		}

		if msg.MessageType == MessageWithPayloadType {
			payloadQuery := `
				SELECT file_path, file_name, content_type, file_size
				FROM public.message_payload
				WHERE message_id = $1
			`
			payloadRows, err := r.db.QueryContext(ctx, payloadQuery, msg.ID)
			if err != nil {
				log.Println("get payloads:", err)
				return nil, ErrDatabaseOperation
			}

			for payloadRows.Next() {
				var path, filename, contentType string
				var size int64
				err := payloadRows.Scan(&path, &filename, &contentType, &size)
				if err != nil {
					payloadRows.Close()
					log.Printf("scan payload error: %v", err)
					return nil, err
				}
				payload := model.Payload{
					URL:         path,
					Filename:    filename,
					Size:        size,
					ContentType: contentType,
				}
				if contentType == filePayloadType {
					msg.FilesDTO = append(msg.FilesDTO, payload)
				} else if contentType == photoPayloadType {
					msg.PhotosDTO = append(msg.PhotosDTO, payload)
				}
			}
			payloadRows.Close()
		}

		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		log.Println("rows error:", err)
		return nil, ErrDatabaseOperation
	}

	return messages, nil
}

func (r *messageRepo) GetMessagesBefore(ctx context.Context, chatID, beforeMessageID uuid.UUID) ([]model.Message, error) {
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
			u.username,
			m.message_type,
			m.sticker_path
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

	rows, err := r.db.QueryContext(ctx, query, chatID, beforeMessageID, 5)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString
		var stickerPath sql.NullString

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
			&msg.MessageType,
			&stickerPath,
		); err != nil {
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			id, err := uuid.Parse(parentMsgID.String)
			if err == nil {
				msg.ParentMessageID = &id
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		if stickerPath.Valid {
			msg.Sticker = stickerPath.String
		}

		if msg.MessageType == MessageWithPayloadType {
			payloadQuery := `
				SELECT file_path, file_name, content_type, file_size
				FROM public.message_payload
				WHERE message_id = $1
			`
			payloadRows, err := r.db.QueryContext(ctx, payloadQuery, msg.ID)
			if err != nil {
				log.Println("get payloads:", err)
				return nil, ErrDatabaseOperation
			}

			for payloadRows.Next() {
				var path, filename, contentType string
				var size int64
				err := payloadRows.Scan(&path, &filename, &contentType, &size)
				if err != nil {
					payloadRows.Close()
					log.Printf("scan payload error: %v", err)
					return nil, err
				}
				payload := model.Payload{
					URL:         path,
					Filename:    filename,
					Size:        size,
					ContentType: contentType,
				}
				if contentType == filePayloadType {
					msg.FilesDTO = append(msg.FilesDTO, payload)
				} else if contentType == photoPayloadType {
					msg.PhotosDTO = append(msg.PhotosDTO, payload)
				}
			}
			payloadRows.Close()
		}

		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		log.Println("rows error:", err)
		return nil, ErrDatabaseOperation
	}

	return messages, nil
}

func (r *messageRepo) GetMessagesAfter(ctx context.Context, chatID, afterMessageID uuid.UUID) ([]model.Message, error) {
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
			u.username,
			m.message_type,
			m.sticker_path
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

	rows, err := r.db.QueryContext(ctx, query, chatID, afterMessageID, 5)
	if err != nil {
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		var parentMsgID sql.NullString
		var avatarPath sql.NullString
		var stickerPath sql.NullString

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
			&msg.MessageType,
			&stickerPath,
		); err != nil {
			return nil, ErrDatabaseScan
		}

		if parentMsgID.Valid {
			id, err := uuid.Parse(parentMsgID.String)
			if err == nil {
				msg.ParentMessageID = &id
			}
		}

		if avatarPath.Valid {
			msg.AvatarPath = &avatarPath.String
		}

		if stickerPath.Valid {
			msg.Sticker = stickerPath.String
		}

		if msg.MessageType == MessageWithPayloadType {
			payloadQuery := `
				SELECT file_path, file_name, content_type, file_size
				FROM public.message_payload
				WHERE message_id = $1
			`
			payloadRows, err := r.db.QueryContext(ctx, payloadQuery, msg.ID)
			if err != nil {
				log.Println("get payloads:", err)
				return nil, ErrDatabaseOperation
			}

			for payloadRows.Next() {
				var path, filename, contentType string
				var size int64
				err := payloadRows.Scan(&path, &filename, &contentType, &size)
				if err != nil {
					payloadRows.Close()
					log.Printf("scan payload error: %v", err)
					return nil, err
				}
				payload := model.Payload{
					URL:         path,
					Filename:    filename,
					Size:        size,
					ContentType: contentType,
				}
				if contentType == filePayloadType {
					msg.FilesDTO = append(msg.FilesDTO, payload)
				} else if contentType == photoPayloadType {
					msg.PhotosDTO = append(msg.PhotosDTO, payload)
				}
			}
			payloadRows.Close()
		}

		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		log.Println("rows error:", err)
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
			u.avatar_path,
			u.username,
			m.message_type,
			m.sticker_path
		FROM message m
		JOIN public.user u ON m.user_id = u.id
		WHERE m.id = $1
	`

	var msg model.Message
	var parentMsgID sql.NullString
	var avatarPath sql.NullString
	var messageType string
	var stickerPath sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&msg.ID,
		&parentMsgID,
		&msg.ChatID,
		&msg.UserID,
		&msg.Body,
		&msg.SentAt,
		&msg.IsRedacted,
		&avatarPath,
		&msg.Username,
		&messageType,
		&stickerPath,
	)
	if err != nil {
		log.Println("get message:", err)
		return nil, ErrDatabaseOperation
	}

	msg.MessageType = messageType

	if parentMsgID.Valid {
		parentID, err := uuid.Parse(parentMsgID.String)
		if err == nil {
			msg.ParentMessageID = &parentID
		}
	}

	if avatarPath.Valid {
		msg.AvatarPath = &avatarPath.String
	}

	if stickerPath.Valid {
		msg.Sticker = stickerPath.String
	}

	if messageType == MessageWithPayloadType {
		payloadQuery := `
			SELECT file_path, file_name, content_type, file_size
			FROM message_payload
			WHERE message_id = $1
		`
		rows, err := r.db.QueryContext(ctx, payloadQuery, msg.ID)
		if err != nil {
			log.Println("get payload:", err)
			return nil, ErrDatabaseOperation
		}
		defer rows.Close()

		for rows.Next() {
			var path, filename, contentType string
			var size int64

			err = rows.Scan(&path, &filename, &contentType, &size)
			if err != nil {
				log.Printf("scan payload error: %v", err)
				return nil, err
			}

			payload := model.Payload{
				URL:         path,
				Filename:    filename,
				Size:        size,
				ContentType: contentType,
			}

			switch contentType {
			case filePayloadType:
				msg.FilesDTO = append(msg.FilesDTO, payload)
			case photoPayloadType:
				msg.PhotosDTO = append(msg.PhotosDTO, payload)
			}
		}
	}

	return &msg, nil
}

func (r *messageRepo) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	messageType := defaultMessageType
	if message.Sticker != "" {
		messageType = stickerMessageType
	} else if len(message.Files) > 0 || len(message.Photos) > 0 {
		messageType = MessageWithPayloadType
	}

	query := `
		INSERT INTO message (user_id, chat_id, body, message_type, sticker_path)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	log.Println("CreateMessage chatID:", message.ChatID)

	err := r.db.QueryRowContext(ctx, query,
		message.UserID,
		message.ChatID,
		message.Body,
		messageType,
		message.Sticker,
	).Scan(&message.ID)
	if err != nil {
		log.Println("insert message:", err)
		return nil, ErrDatabaseOperation
	}

	if messageType == MessageWithPayloadType {
		for _, file := range message.FilesDTO {
			id := uuid.New()
			_, err = r.db.ExecContext(ctx, `
				INSERT INTO message_payload (id, message_id, file_path, file_name, content_type, file_size)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, id, message.ID, file.URL, file.Filename, file.ContentType, file.Size)
			if err != nil {
				log.Println("insert file payload:", err)
				return nil, ErrDatabaseOperation
			}
		}

		for _, photo := range message.PhotosDTO {
			id := uuid.New()
			_, err = r.db.ExecContext(ctx, `
				INSERT INTO message_payload (id, message_id, file_path, file_name, content_type, file_size)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, id, message.ID, photo.URL, photo.Filename, photo.ContentType, photo.Size)
			if err != nil {
				log.Println("insert photo payload:", err)
				return nil, ErrDatabaseOperation
			}
		}
	}

	messageOut, err := r.GetMessage(ctx, message.ID)
	if err != nil {
		log.Println("get message:", err)
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
