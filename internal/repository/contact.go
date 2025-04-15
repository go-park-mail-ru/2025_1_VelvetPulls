package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IContactRepo interface {
	GetContacts(ctx context.Context, userID uuid.UUID) ([]model.Contact, error)
	AddContactByUsername(ctx context.Context, userID uuid.UUID, contactUsername string) error
	DeleteContactByUsername(ctx context.Context, userID uuid.UUID, contactUsername string) error
}

type contactRepo struct {
	db *sql.DB
}

func NewContactRepo(db *sql.DB) IContactRepo {
	return &contactRepo{db: db}
}

func (r *contactRepo) GetContacts(ctx context.Context, userID uuid.UUID) ([]model.Contact, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Fetching contacts", zap.String("userID", userID.String()))

	query := `
		SELECT u.id, u.first_name, u.last_name, u.username, u.avatar_path
		FROM public.contact c
		JOIN public.user u ON c.contact_id = u.id
		WHERE c.user_id = $1`

	logger.Debug("Executing query", zap.String("query", query), zap.String("userID", userID.String()))
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.Error("Failed to fetch contacts",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var contacts []model.Contact
	for rows.Next() {
		var contact model.Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Username, &contact.AvatarURL); err != nil {
			logger.Error("Failed to scan contact",
				zap.String("userID", userID.String()),
				zap.Error(err),
			)
			return nil, ErrDatabaseScan
		}
		logger.Debug("Fetched contact", zap.Any("contact", contact))
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error iterating over contacts",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return nil, ErrDatabaseOperation
	}

	logger.Info("Successfully fetched contacts",
		zap.String("userID", userID.String()),
		zap.Int("count", len(contacts)),
	)
	return contacts, nil
}

func (r *contactRepo) AddContactByUsername(ctx context.Context, userID uuid.UUID, contactUsername string) error {
	logger := utils.GetLoggerFromCtx(ctx)

	cleanUsername := utils.SanitizeString(contactUsername)
	if cleanUsername == "" {
		logger.Warn("Empty username after sanitization")
		return ErrEmptyField
	}

	logger.Info("Adding contact by username",
		zap.String("userID", userID.String()),
		zap.String("contactUsername", cleanUsername),
	)

	var contactID uuid.UUID
	err := r.db.QueryRowContext(ctx,
		"SELECT id FROM public.user WHERE username = $1",
		cleanUsername,
	).Scan(&contactID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Contact user not found", zap.String("username", contactUsername))
			return ErrUserNotFound
		}
		logger.Error("Failed to fetch contact ID", zap.Error(err))
		return ErrDatabaseOperation
	}

	// Самозапрещаем: нельзя добавить самого себя
	if userID == contactID {
		logger.Warn("User attempted to add themselves as a contact")
		return ErrSelfContact
	}

	query := `INSERT INTO public.contact (user_id, contact_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	logger.Debug("Executing insert",
		zap.String("query", query),
		zap.String("userID", userID.String()),
		zap.String("contactID", contactID.String()),
	)

	res, err := r.db.ExecContext(ctx, query, userID, contactID)
	if err != nil {
		logger.Error("Failed to insert contact", zap.Error(err))
		return ErrDatabaseOperation
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("Failed to check affected rows", zap.Error(err))
		return ErrDatabaseOperation
	}

	if rowsAffected == 0 {
		logger.Info("Contact already exists")
	} else {
		logger.Info("Contact successfully added")
	}

	return nil
}

func (r *contactRepo) DeleteContactByUsername(ctx context.Context, userID uuid.UUID, contactUsername string) error {
	logger := utils.GetLoggerFromCtx(ctx)

	cleanUsername := utils.SanitizeString(contactUsername)
	if cleanUsername == "" {
		logger.Warn("Empty username after sanitization")
		return ErrEmptyField
	}

	logger.Info("Deleting contact by username",
		zap.String("userID", userID.String()),
		zap.String("contactUsername", cleanUsername),
	)

	var contactID uuid.UUID
	err := r.db.QueryRowContext(ctx,
		"SELECT id FROM public.user WHERE username = $1",
		cleanUsername,
	).Scan(&contactID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Contact user not found", zap.String("username", cleanUsername))
			return ErrUserNotFound
		}
		logger.Error("Failed to fetch contact ID", zap.Error(err))
		return ErrDatabaseOperation
	}

	query := `DELETE FROM public.contact WHERE user_id = $1 AND contact_id = $2`
	logger.Debug("Executing delete",
		zap.String("query", query),
		zap.String("userID", userID.String()),
		zap.String("contactID", contactID.String()),
	)

	res, err := r.db.ExecContext(ctx, query, userID, contactID)
	if err != nil {
		logger.Error("Failed to delete contact", zap.Error(err))
		return ErrDatabaseOperation
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("Failed to check affected rows", zap.Error(err))
		return ErrDatabaseOperation
	}

	if rowsAffected == 0 {
		logger.Warn("Contact not found for deletion",
			zap.String("userID", userID.String()),
			zap.String("contactUsername", cleanUsername),
		)
	} else {
		logger.Info("Successfully deleted contact",
			zap.String("userID", userID.String()),
			zap.String("contactUsername", cleanUsername),
		)
	}

	return nil
}
