package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

type IAuthRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type authRepo struct {
	db *sql.DB
}

func NewauthRepo(db *sql.DB) IAuthRepo {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
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
	FROM users 
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
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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
	FROM users 
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
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
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
	FROM users 
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
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepo) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users 
	(
	username,
	phone, 
	password
	) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, user.Username, user.Phone, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
