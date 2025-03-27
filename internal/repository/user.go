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
	GetAvatarPathByUserID(ctx context.Context, userID string) (string, error)
	UpdateAvatarPathByUserID(ctx context.Context, userID string, avatarPath string) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (string, error)
	UpdateUser(ctx context.Context, user *model.User) error
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
		fmt.Print(err)
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

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
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

func (r *userRepo) UpdateUser(ctx context.Context, user *model.User) error {
	var (
		queryBuilder strings.Builder
		args         []interface{}
		argIndex     = 1
	)

	queryBuilder.WriteString("UPDATE public.user SET ")

	if user.FirstName != nil {
		queryBuilder.WriteString(fmt.Sprintf("first_name = $%d, ", argIndex))
		args = append(args, user.FirstName)
		argIndex++
	}
	if user.LastName != nil {
		queryBuilder.WriteString(fmt.Sprintf("last_name = $%d, ", argIndex))
		args = append(args, user.LastName)
		argIndex++
	}
	if user.Username != "" {
		queryBuilder.WriteString(fmt.Sprintf("username = $%d, ", argIndex))
		args = append(args, user.Username)
		argIndex++
	}
	if user.Phone != "" {
		queryBuilder.WriteString(fmt.Sprintf("phone = $%d, ", argIndex))
		args = append(args, user.Phone)
		argIndex++
	}
	if user.Email != nil {
		queryBuilder.WriteString(fmt.Sprintf("email = $%d, ", argIndex))
		args = append(args, user.Email)
		argIndex++
	}

	queryBuilder.WriteString(fmt.Sprintf("updated_at = $%d ", argIndex))
	args = append(args, time.Now())
	argIndex++

	queryBuilder.WriteString(fmt.Sprintf("WHERE id = $%d", argIndex))
	args = append(args, user.ID)

	_, err := r.db.ExecContext(ctx, queryBuilder.String(), args...)
	return err
}

func (r *userRepo) GetAvatarPathByUserID(ctx context.Context, userID string) (string, error) {
	var avatarPath sql.NullString
	query := `SELECT avatar_path 
	FROM public.user 
	WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&avatarPath)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	if !avatarPath.Valid {
		return "", ErrAvatarNotFound
	}

	return avatarPath.String, nil
}

func (r *userRepo) UpdateAvatarPathByUserID(ctx context.Context, userID string, avatarPath string) error {

	query := `UPDATE public.user 
	SET 
	avatar_path = $1, 
	updated_at = $2 
	WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, avatarPath, time.Now(), userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
