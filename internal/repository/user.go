package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IUserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UpdateUserProfile) (string, string, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) IUserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) getUserByField(ctx context.Context, field string, value interface{}) (*model.User, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	if value == "" || value == uuid.Nil {
		logger.Warn("empty lookup value", zap.String("field", field))
		return nil, ErrEmptyField
	}

	query := fmt.Sprintf(
		`SELECT id, avatar_path, name, username, password, birth_date FROM public.user WHERE %s = $1`, field,
	)

	logger.Info("executing user lookup", zap.String("field", field))

	var u model.User
	err := r.db.QueryRowContext(ctx, query, value).Scan(
		&u.ID, &u.AvatarPath, &u.Name,
		&u.Username, &u.Password, &u.BirthDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("user not found", zap.String("field", field))
			return nil, ErrUserNotFound
		}
		logger.Error("lookup failed", zap.Error(err))
		return nil, ErrDatabaseOperation
	}
	return &u, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.getUserByField(ctx, "username", username)
}

func (r *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if id == uuid.Nil {
		utils.GetLoggerFromCtx(ctx).Warn("invalid uuid provided")
		return nil, ErrInvalidUUID
	}
	return r.getUserByField(ctx, "id", id.String())
}

func (r *userRepo) UpdateUser(ctx context.Context, profile *model.UpdateUserProfile) (string, string, error) {
	if profile == nil {
		utils.GetLoggerFromCtx(ctx).Warn("empty profile data")
		return "", "", ErrEmptyField
	}
	if profile.ID == uuid.Nil {
		utils.GetLoggerFromCtx(ctx).Warn("invalid uuid provided")
		return "", "", ErrInvalidUUID
	}

	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("updating user", zap.String("userID", profile.ID.String()))

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("begin tx failed", zap.Error(err))
		return "", "", ErrDatabaseOperation
	}

	var (
		updates              []string
		args                 []interface{}
		idx                  = 1
		oldAvatar, newAvatar string
	)

	if profile.Avatar != nil {
		logger.Info("updating avatar")
		var current *string
		err = tx.QueryRowContext(ctx, "SELECT avatar_path FROM public.user WHERE id = $1 FOR UPDATE", profile.ID).Scan(&current)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			rollbackTx(logger, tx)
			return "", "", ErrDatabaseOperation
		}
		if current != nil {
			oldAvatar = *current
		}
		newAvatar = fmt.Sprintf("/uploads/avatar/%s.png", uuid.New().String())
		updates = append(updates, fmt.Sprintf("avatar_path = $%d", idx))
		args = append(args, newAvatar)
		idx++
	}

	if profile.Password != "" {
		logger.Info("updating password")
		hashedPassword, err := utils.HashAndSalt(profile.Password)
		if err != nil {
			rollbackTx(logger, tx)
			return "", "", ErrDatabaseOperation
		}
		updates = append(updates, fmt.Sprintf("password = $%d", idx))
		args = append(args, hashedPassword)
		idx++
	}

	fields := map[string]*string{
		"name":     profile.Name,
		"username": profile.Username,
	}
	for field, ptr := range fields {
		if ptr != nil {
			if *ptr == "" {
				rollbackTx(logger, tx)
				return "", "", ErrEmptyField
			}
			updates = append(updates, fmt.Sprintf("%s = $%d", field, idx))
			args = append(args, *ptr)
			idx++
		}
	}

	if profile.BirthDate != nil {
		updates = append(updates, fmt.Sprintf("birth_date = $%d", idx))
		args = append(args, *profile.BirthDate)
		idx++
	}

	if len(updates) == 0 {
		rollbackTx(logger, tx)
		return "", "", nil
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", idx))
	args = append(args, time.Now())
	idx++

	args = append(args, profile.ID)
	query := fmt.Sprintf("UPDATE public.user SET %s WHERE id = $%d", strings.Join(updates, ", "), idx)

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("duplicate update attempt", zap.Error(err))
			rollbackTx(logger, tx)
			return "", "", ErrRecordAlreadyExists
		}
		logger.Error("update failed", zap.Error(err))
		rollbackTx(logger, tx)
		return "", "", ErrDatabaseOperation
	}

	if n, _ := res.RowsAffected(); n == 0 {
		logger.Warn("no rows affected on update")
		rollbackTx(logger, tx)
		return "", "", ErrUpdateFailed
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit failed", zap.Error(err))
		return "", "", ErrDatabaseOperation
	}

	return newAvatar, oldAvatar, nil
}
