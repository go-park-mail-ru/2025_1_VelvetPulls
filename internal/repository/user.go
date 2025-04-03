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

	var updates []string
	var args []interface{}
	argIndex := 1
	var avatarOldURL, avatarNewURL string

	// Обработка нового аватара
	if profile.Avatar != nil {
		logger.Info("Updating avatar")
		var oldUrl *string
		err := r.db.QueryRowContext(ctx, "SELECT avatar_path FROM public.user WHERE id = $1", profile.ID).Scan(&oldUrl)
		if oldUrl != nil {
			avatarOldURL = *oldUrl
		}
		fmt.Println(err)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("Failed to get current avatar path")
			return "", "", ErrDatabaseOperation
		}

		// Формируем путь для нового аватара
		avatarDir := "./uploads/avatar/"
		avatarNewURL = avatarDir + uuid.New().String() + ".png"
		updates = append(updates, fmt.Sprintf("avatar_path = $%d", argIndex))
		args = append(args, avatarNewURL)
		argIndex++
	}

	// Обработка остальных полей
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
				return "", "", ErrEmptyField
			}
			updates = append(updates, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, *value)
			argIndex++
		}
	}

	// Если нет обновлений, возвращаем пустые строки
	if len(updates) == 0 {
		return "", "", nil
	}

	// Добавляем время обновления
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Строим запрос на обновление
	query := fmt.Sprintf("UPDATE public.user SET %s WHERE id = $%d", strings.Join(updates, ", "), argIndex)
	args = append(args, profile.ID)

	// Выполняем запрос
	logger.Info("Executing update query")
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			logger.Warn("Attempt to update to existing values")
			return "", "", ErrRecordAlreadyExists
		}
		logger.Error("Database operation failed")
		return "", "", ErrDatabaseOperation
	}

	// Проверяем количество измененных строк
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("Failed to get affected rows count")
		return "", "", ErrDatabaseOperation
	}
	if rowsAffected == 0 {
		logger.Warn("No rows affected, update failed")
		return "", "", ErrUpdateFailed
	}

	return avatarNewURL, avatarOldURL, nil
}
