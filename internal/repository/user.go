package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
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
	if value == "" {
		return nil, ErrEmptyField
	}

	var user model.User
	query := fmt.Sprintf(`SELECT id, first_name, last_name, username, phone, email, password, created_at, updated_at FROM public.user WHERE %s = $1`, field)
	row := r.db.QueryRowContext(ctx, query, value)

	err := row.Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Phone, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabaseOperation
	}
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
		return nil, ErrInvalidUUID
	}
	return r.getUserByField(ctx, "id", id.String())
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (string, error) {
	if user == nil {
		return "", ErrEmptyField
	}
	if user.Username == "" || user.Phone == "" || user.Password == "" {
		return "", ErrEmptyField
	}

	query := `INSERT INTO public.user (username, phone, password) VALUES ($1, $2, $3) RETURNING id`
	var userID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Phone, user.Password).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "", ErrRecordAlreadyExists
		}
		return "", ErrDatabaseOperation
	}
	return userID.String(), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, profile *model.UpdateUserProfile) (string, string, error) {
	if profile == nil {
		return "", "", ErrEmptyField
	}
	if profile.ID == uuid.Nil {
		return "", "", ErrInvalidUUID
	}

	var updates []string
	var args []interface{}
	argIndex := 1
	var avatarOldURL, avatarNewURL string

	// Обработка нового аватара
	if profile.Avatar != nil {
		err := r.db.QueryRowContext(ctx, "SELECT avatar_path FROM public.user WHERE id = $1", profile.ID).Scan(&avatarOldURL)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", "", ErrDatabaseOperation
		}

		currentDate := time.Now()
		year := currentDate.Year()
		month := currentDate.Month()
		day := currentDate.Day()

		// Формируем путь для нового аватара
		avatarDir := fmt.Sprintf("./uploads/avatar/%d/%02d/%02d/", year, month, day)
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
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "", "", ErrRecordAlreadyExists
		}
		return "", "", ErrDatabaseOperation
	}

	// Проверяем количество измененных строк
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", "", ErrDatabaseOperation
	}
	if rowsAffected == 0 {
		return "", "", ErrUpdateFailed
	}

	return avatarNewURL, avatarOldURL, nil
}
