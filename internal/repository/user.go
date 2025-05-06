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

// IUserRepo описывает операции для работы с пользователями
type IUserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UpdateUserProfile) (string, string, error)
}

// userRepo реализует IUserRepo через PostgreSQL
type userRepo struct {
	db *sql.DB
}

// NewUserRepo создаёт новый репозиторий пользователей
func NewUserRepo(db *sql.DB) IUserRepo {
	return &userRepo{db: db}
}

// getUserByField объединяет логику поиска пользователя по любому полю
func (r *userRepo) getUserByField(ctx context.Context, field string, value interface{}) (*model.User, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	if value == "" || value == uuid.Nil {
		logger.Warn("empty lookup value", zap.String("field", field))
		return nil, ErrEmptyField
	}

	query := fmt.Sprintf(
		`SELECT id, avatar_path, first_name, last_name, username, phone, email, password, created_at, updated_at FROM public.user WHERE %s = $1`,
		field,
	)
	logger.Info("executing user lookup", zap.String("field", field))

	var u model.User
	err := r.db.QueryRowContext(ctx, query, value).Scan(
		&u.ID, &u.AvatarPath, &u.FirstName, &u.LastName,
		&u.Username, &u.Phone, &u.Email, &u.Password,
		&u.CreatedAt, &u.UpdatedAt,
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

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("creating user")

	if user == nil || user.Username == "" || user.Phone == "" || user.Password == "" {
		logger.Warn("missing required fields for create")
		return uuid.Nil, ErrEmptyField
	}

	query := `INSERT INTO public.user (username, phone, password) VALUES ($1, $2, $3) RETURNING id`
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Phone, user.Password).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("user exists", zap.Error(err))
			return uuid.Nil, ErrRecordAlreadyExists
		}
		logger.Error("create failed", zap.Error(err))
		return uuid.Nil, ErrDatabaseOperation
	}
	return id, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, profile *model.UpdateUserProfile) (string, string, error) {
	// проверяем входные параметры
	if profile == nil {
		utils.GetLoggerFromCtx(ctx).Warn("empty profile data")
		return "", "", ErrEmptyField
	}
	if profile.ID == uuid.Nil {
		utils.GetLoggerFromCtx(ctx).Warn("invalid uuid provided")
		return "", "", ErrInvalidUUID
	}

	// теперь безопасно получаем logger и логируем ID
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

	// обновление аватара
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
		newAvatar = fmt.Sprintf("./uploads/avatar/%s.png", uuid.New().String())
		updates = append(updates, fmt.Sprintf("avatar_path = $%d", idx))
		args = append(args, newAvatar)
		idx++
	}

	// обновление пароля
	if profile.Password != "" {
		logger.Info("updating password")
		// Хешируем новый пароль
		hashedPassword, err := utils.HashAndSalt(profile.Password)
		if err != nil {
			rollbackTx(logger, tx)
			return "", "", ErrDatabaseOperation
		}

		updates = append(updates, fmt.Sprintf("password = $%d", idx))
		args = append(args, hashedPassword)
		idx++
	}

	// обновление полей
	fields := map[string]*string{
		"first_name": profile.FirstName,
		"last_name":  profile.LastName,
		"username":   profile.Username,
		"phone":      profile.Phone,
		"email":      profile.Email,
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

	if len(updates) == 0 {
		rollbackTx(logger, tx)
		return "", "", nil
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", idx))
	args = append(args, time.Now())
	idx++

	// финальный запрос
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
