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
	UpdateUser(ctx context.Context, user *model.UpdateUserProfile) (*string, *string, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) IUserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query := `SELECT 
	id, 
	first_name, 
	last_name, 
	username, 
	phone, 
	email, 
	password, 
	created_at, 
	updated_at 
	FROM public.user 
	WHERE username = $1`
	row := r.db.QueryRow(query, username)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT 
	id, 
	first_name, 
	last_name, 
	username, 
	phone, 
	email, 
	password, 
	created_at, 
	updated_at 
	FROM public.user 
	WHERE email = $1`
	row := r.db.QueryRow(query, email)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	query := `SELECT 
	id, 
	first_name, 
	last_name, 
	username, 
	phone, 
	email,
	password, 
	created_at, 
	updated_at 
	FROM public.user
	WHERE phone = $1`
	row := r.db.QueryRow(query, phone)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	query := `SELECT 
	id, 
	avatar_path,
	first_name, 
	last_name, 
	username, 
	phone, 
	email,
	password, 
	created_at, 
	updated_at 
	FROM public.user
	WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.AvatarPath,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (string, error) {
	query := `INSERT INTO public.user 
	(
	username,
	phone, 
	password
	) VALUES ($1, $2, $3) RETURNING id`
	var userID uuid.UUID
	err := r.db.QueryRow(query, user.Username, user.Phone, user.Password).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID.String(), nil
}

func (r *userRepo) UpdateUser(ctx context.Context, profile *model.UpdateUserProfile) (*string, *string, error) {
	var (
		queryBuilder strings.Builder
		args         []interface{}
		argIndex     = 1
		avatarOldURL *string
		avatarNewURL *string
		updates      []string
	)

	// Получаем старый путь аватарки если нужно
	if profile.Avatar != nil {
		err := r.db.QueryRowContext(ctx,
			"SELECT avatar_path FROM public.user WHERE id = $1",
			profile.ID,
		).Scan(&avatarOldURL)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, fmt.Errorf("failed to get old avatar path: %w", err)
		}
	}

	// Формируем части запроса
	if profile.FirstName != nil {
		updates = append(updates, fmt.Sprintf("first_name = $%d", argIndex))
		args = append(args, *profile.FirstName)
		argIndex++
	}
	if profile.LastName != nil {
		updates = append(updates, fmt.Sprintf("last_name = $%d", argIndex))
		args = append(args, *profile.LastName)
		argIndex++
	}
	if profile.Username != nil {
		updates = append(updates, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, *profile.Username)
		argIndex++
	}
	if profile.Phone != nil {
		updates = append(updates, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, *profile.Phone)
		argIndex++
	}
	if profile.Email != nil {
		updates = append(updates, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *profile.Email)
		argIndex++
	}
	if profile.Avatar != nil {
		updates = append(updates, fmt.Sprintf("avatar_path = $%d", argIndex))
		avatarNewURL = new(string)
		*avatarNewURL = "./uploads/avatar/" + uuid.New().String() + ".png"
		args = append(args, *avatarNewURL)
		argIndex++
	}

	// Если нет полей для обновления
	if len(updates) == 0 {
		return nil, nil, nil
	}

	// Добавляем updated_at
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Формируем итоговый запрос
	queryBuilder.WriteString("UPDATE public.user SET ")
	queryBuilder.WriteString(strings.Join(updates, ", "))
	queryBuilder.WriteString(fmt.Sprintf(" WHERE id = $%d", argIndex))
	args = append(args, profile.ID)

	_, err := r.db.ExecContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, nil, fmt.Errorf("update failed: %w", err)
	}

	return avatarNewURL, avatarOldURL, nil
}
