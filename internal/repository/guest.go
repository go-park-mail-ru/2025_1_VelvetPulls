package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IGuestMapper interface {
	CreateGuest(ctx context.Context, guest *model.Guest) (uuid.UUID, error)
	GetGuestByID(ctx context.Context, id uuid.UUID) (*model.Guest, error)
	GetGuestByTemporaryID(ctx context.Context, tempID string) (*model.Guest, error)
	UpdateGuest(ctx context.Context, guest *model.Guest) error
	DeleteGuest(ctx context.Context, id uuid.UUID) error
	GetGuestBySession(ctx context.Context, sessionID string) (*model.Guest, error)
}

type GuestMapper struct {
	db *sql.DB
}

func NewGuestMapper(db *sql.DB) IGuestMapper {
	return &GuestMapper{db: db}
}

func (r *GuestMapper) CreateGuest(ctx context.Context, guest *model.Guest) (uuid.UUID, error) {
	logger := zap.L()
	
	query := `
		INSERT INTO guests (temporary_id, status, created_at, last_activity_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id
	`
	
	var guestID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, 
		guest.TemporaryID,
		guest.Status,
	).Scan(&guestID)
	
	if err != nil {
		logger.Error("CreateGuest: failed to create guest", zap.Error(err))
		return uuid.Nil, err
	}
	
	logger.Info("CreateGuest: guest created successfully", zap.String("guestID", guestID.String()))
	return guestID, nil
}

func (r *GuestMapper) GetGuestByID(ctx context.Context, id uuid.UUID) (*model.Guest, error) {
	logger := zap.L()
	
	query := `
		SELECT id, temporary_id, status, created_at, last_activity_at, session_id
		FROM guests
		WHERE id = $1 AND deleted_at IS NULL
	`
	
	var guest model.Guest
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&guest.ID,
		&guest.TemporaryID,
		&guest.Status,
		&guest.CreatedAt,
		&guest.LastActivityAt,
		&guest.SessionID,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("GetGuestByID: guest not found", zap.String("guestID", id.String()))
			return nil, ErrGuestNotFound
		}
		logger.Error("GetGuestByID: failed to get guest", zap.Error(err))
		return nil, err
	}
	
	return &guest, nil
}

func (r *GuestMapper) GetGuestByTemporaryID(ctx context.Context, tempID string) (*model.Guest, error) {
	logger := zap.L()
	
	query := `
		SELECT id, temporary_id, status, created_at, last_activity_at, session_id
		FROM guests
		WHERE temporary_id = $1 AND deleted_at IS NULL
	`
	
	var guest model.Guest
	err := r.db.QueryRowContext(ctx, query, tempID).Scan(
		&guest.ID,
		&guest.TemporaryID,
		&guest.Status,
		&guest.CreatedAt,
		&guest.LastActivityAt,
		&guest.SessionID,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("GetGuestByTemporaryID: guest not found", zap.String("tempID", tempID))
			return nil, ErrGuestNotFound
		}
		logger.Error("GetGuestByTemporaryID: failed to get guest", zap.Error(err))
		return nil, err
	}
	
	return &guest, nil
}

func (r *GuestMapper) GetGuestBySession(ctx context.Context, sessionID string) (*model.Guest, error) {
	logger := zap.L()
	
	query := `
		SELECT id, temporary_id, status, created_at, last_activity_at, session_id
		FROM guests
		WHERE session_id = $1 AND deleted_at IS NULL
	`
	
	var guest model.Guest
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&guest.ID,
		&guest.TemporaryID,
		&guest.Status,
		&guest.CreatedAt,
		&guest.LastActivityAt,
		&guest.SessionID,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("GetGuestBySession: guest not found", zap.String("sessionID", sessionID))
			return nil, ErrGuestNotFound
		}
		logger.Error("GetGuestBySession: failed to get guest", zap.Error(err))
		return nil, err
	}
	
	return &guest, nil
}

func (r *GuestMapper) UpdateGuest(ctx context.Context, guest *model.Guest) error {
	logger := zap.L()
	
	query := `
		UPDATE guests
		SET temporary_id = $1, 
		    status = $2, 
		    last_activity_at = NOW(),
		    session_id = $3,
		    updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`
	
	result, err := r.db.ExecContext(ctx, query,
		guest.TemporaryID,
		guest.Status,
		guest.SessionID,
		guest.ID,
	)
	
	if err != nil {
		logger.Error("UpdateGuest: failed to update guest", 
			zap.String("guestID", guest.ID.String()), 
			zap.Error(err))
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		logger.Warn("UpdateGuest: no guest found to update", zap.String("guestID", guest.ID.String()))
		return ErrGuestNotFound
	}
	
	logger.Info("UpdateGuest: guest updated successfully", zap.String("guestID", guest.ID.String()))
	return nil
}

func (r *GuestMapper) DeleteGuest(ctx context.Context, id uuid.UUID) error {
	logger := zap.L()
	
	query := `
		UPDATE guests
		SET deleted_at = NOW(), status = 'deleted'
		WHERE id = $1 AND deleted_at IS NULL
	`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("DeleteGuest: failed to delete guest", 
			zap.String("guestID", id.String()), 
			zap.Error(err))
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		logger.Warn("DeleteGuest: no guest found to delete", zap.String("guestID", id.String()))
		return ErrGuestNotFound
	}
	
	logger.Info("DeleteGuest: guest deleted successfully", zap.String("guestID", id.String()))
	return nil
}