package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IAuthRepo interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) IAuthRepo {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	if user == nil {
		logger.Warn("User data is nil")
		return uuid.Nil, ErrEmptyField
	}
	if user.Username == "" || user.Password == "" {
		logger.Warn("Missing required fields", zap.String("username", user.Username))
		return uuid.Nil, ErrEmptyField
	}

	// Убираем поле birth_date из запроса
	query := `INSERT INTO public.user (username, password, name) 
              VALUES ($1, $2, $3) RETURNING id`
	var userID uuid.UUID

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Name).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("User already exists", zap.Error(err))
			return uuid.Nil, ErrRecordAlreadyExists
		}
		logger.Error("CreateUser query failed", zap.Error(err))
		return uuid.Nil, ErrDatabaseOperation
	}

	logger.Info("User created", zap.String("user_id", userID.String()))
	return userID, nil
}

func (r *authRepo) getUserByField(ctx context.Context, field, value string) (*model.User, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	if value == "" {
		logger.Warn("Field value is empty", zap.String("field", field))
		return nil, ErrEmptyField
	}

	var user model.User
	// Обновляем запрос, чтобы не учитывать поле birth_date при извлечении пользователя
	query := fmt.Sprintf(`SELECT id, avatar_path, name, username
                          FROM public.user WHERE %s = $1`, field)
	logger.Info("Executing query to get user by field")
	row := r.db.QueryRowContext(ctx, query, value)

	err := row.Scan(
		&user.ID, &user.AvatarPath, &user.Name, &user.Username,
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

func (r *authRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.getUserByField(ctx, "username", username)
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.getUserByField(ctx, "email", email)
}

func (r *authRepo) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	return r.getUserByField(ctx, "phone", phone)
}

func (r *authRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return r.getUserByField(ctx, "id", id.String())
}
