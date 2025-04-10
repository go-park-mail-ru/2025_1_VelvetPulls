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
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (string, error)
	UpdateUser(ctx context.Context, user *model.UpdateUserProfile) (string, string, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) IUserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) getUserByField(ctx context.Context, field, value string) (*model.User, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	if value == "" {
		logger.Warn("Field value is empty", zap.String("field", field))
		return nil, ErrEmptyField
	}

	var user model.User
	query := fmt.Sprintf(`SELECT id, avatar_path, first_name, last_name, username, phone, email, password, created_at, updated_at FROM public.user WHERE %s = $1`, field)
	logger.Info("Executing query to get user by field")
	row := r.db.QueryRowContext(ctx, query, value)

	err := row.Scan(
		&user.ID, &user.AvatarPath, &user.FirstName, &user.LastName, &user.Username, &user.Phone, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("User not found")
			return nil, ErrUserNotFound
		}
		logger.Error("Database operation failed")
		return nil, ErrDatabaseOperation
	}
	logger.Info("User found")
	return &user, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.getUserByField(ctx, "username", username)
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.getUserByField(ctx, "email", email)
}

func (r *userRepo) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	return r.getUserByField(ctx, "phone", phone)
}

func (r *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if id == uuid.Nil {
		logger := utils.GetLoggerFromCtx(ctx)
		logger.Warn("Invalid UUID provided")
		return nil, ErrInvalidUUID
	}
	return r.getUserByField(ctx, "id", id.String())
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (string, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	if user == nil {
		logger.Warn("User data is empty")
		return "", ErrEmptyField
	}
	if user.Username == "" || user.Phone == "" || user.Password == "" {
		logger.Warn("Missing required fields for user creation")
		return "", ErrEmptyField
	}

	query := `INSERT INTO public.user (username, phone, password) VALUES ($1, $2, $3) RETURNING id`
	logger.Info("Executing query to create a user")
	var userID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Phone, user.Password).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("User with the same username or phone already exists")
			return "", ErrRecordAlreadyExists
		}
		logger.Error("Database operation failed")
		return "", ErrDatabaseOperation
	}
	return userID.String(), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, profile *model.UpdateUserProfile) (string, string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	if profile == nil {
		logger.Warn("Profile data is empty")
		return "", "", ErrEmptyField
	}
	if profile.ID == uuid.Nil {
		logger.Warn("Invalid UUID provided")
		return "", "", ErrInvalidUUID
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Failed to start transaction", zap.Error(err))
		return "", "", ErrDatabaseOperation
	}

	var updates []string
	var args []interface{}
	argIndex := 1
	var avatarOldURL, avatarNewURL string

	if profile.Avatar != nil {
		logger.Info("Updating avatar")
		err = tx.QueryRowContext(ctx, "SELECT avatar_path FROM public.user WHERE id = $1 FOR UPDATE", profile.ID).Scan(&avatarOldURL)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("Failed to get current avatar path", zap.Error(err))
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error("Rollback failed", zap.Error(rbErr))
			}
			return "", "", ErrDatabaseOperation
		}

		avatarDir := "./uploads/avatar/"
		avatarNewURL = avatarDir + uuid.New().String() + ".png"
		updates = append(updates, fmt.Sprintf("avatar_path = $%d", argIndex))
		args = append(args, avatarNewURL)
		argIndex++
	}

	fields := map[string]*string{
		"first_name": profile.FirstName,
		"last_name":  profile.LastName,
		"username":   profile.Username,
		"phone":      profile.Phone,
		"email":      profile.Email,
	}

	for field, value := range fields {
		if value != nil {
			if *value == "" {
				logger.Warn("Field value is empty")
				if rbErr := tx.Rollback(); rbErr != nil {
					logger.Error("Rollback failed", zap.Error(rbErr))
				}
				return "", "", ErrEmptyField
			}
			updates = append(updates, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, *value)
			argIndex++
		}
	}

	if len(updates) == 0 {
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.Error("Rollback failed", zap.Error(rbErr))
		}
		return "", "", nil
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	query := fmt.Sprintf("UPDATE public.user SET %s WHERE id = $%d", strings.Join(updates, ", "), argIndex)
	args = append(args, profile.ID)

	logger.Info("Executing update query inside transaction")
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("Attempt to update to existing values", zap.Error(err))
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error("Rollback failed", zap.Error(rbErr))
			}
			return "", "", ErrRecordAlreadyExists
		}
		logger.Error("Database operation failed", zap.Error(err))
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.Error("Rollback failed", zap.Error(rbErr))
		}
		return "", "", ErrDatabaseOperation
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("Failed to get affected rows count", zap.Error(err))
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.Error("Rollback failed", zap.Error(rbErr))
		}
		return "", "", ErrDatabaseOperation
	}
	if rowsAffected == 0 {
		logger.Warn("No rows affected, update failed")
		if rbErr := tx.Rollback(); rbErr != nil {
			logger.Error("Rollback failed", zap.Error(rbErr))
		}
		return "", "", ErrUpdateFailed
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction", zap.Error(err))
		return "", "", ErrDatabaseOperation
	}

	return avatarNewURL, avatarOldURL, nil
}
