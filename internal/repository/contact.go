package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IContactRepo interface {
	GetContacts(ctx context.Context, userID uuid.UUID) ([]model.Contact, error)
	AddContact(ctx context.Context, userID, contactID uuid.UUID) error
	DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error
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

func (r *contactRepo) AddContact(ctx context.Context, userID, contactID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Adding contact",
		zap.String("userID", userID.String()),
		zap.String("contactID", contactID.String()),
	)

	query := `INSERT INTO public.contact (user_id, contact_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	logger.Debug("Executing query",
		zap.String("query", query),
		zap.String("userID", userID.String()),
		zap.String("contactID", contactID.String()),
	)

	res, err := r.db.ExecContext(ctx, query, userID, contactID)
	if err != nil {
		logger.Error("Failed to add contact",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
			zap.Error(err),
		)
		return ErrDatabaseOperation
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("Failed to get affected rows",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
			zap.Error(err),
		)
		return ErrDatabaseOperation
	}

	if rowsAffected == 0 {
		logger.Warn("Contact already exists",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
		)
	} else {
		logger.Info("Successfully added contact",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
		)
	}

	return nil
}

func (r *contactRepo) DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Deleting contact",
		zap.String("userID", userID.String()),
		zap.String("contactID", contactID.String()),
	)

	query := `DELETE FROM public.contact WHERE user_id = $1 AND contact_id = $2`
	logger.Debug("Executing query",
		zap.String("query", query),
		zap.String("userID", userID.String()),
		zap.String("contactID", contactID.String()),
	)

	res, err := r.db.ExecContext(ctx, query, userID, contactID)
	if err != nil {
		logger.Error("Failed to delete contact",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
			zap.Error(err),
		)
		return ErrDatabaseOperation
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("Failed to get affected rows",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
			zap.Error(err),
		)
		return ErrDatabaseOperation
	}

	if rowsAffected == 0 {
		logger.Warn("Contact not found for deletion",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
		)
	} else {
		logger.Info("Successfully deleted contact",
			zap.String("userID", userID.String()),
			zap.String("contactID", contactID.String()),
		)
	}

	return nil
}
