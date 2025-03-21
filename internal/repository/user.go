package repository

import (
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

type IUserRepo interface {
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByPhone(phone string) (*model.User, error)
	CreateUser(user *model.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) IUserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	query := "SELECT id, first_name, last_name, username, phone, email, password, created_at, updated_at FROM users WHERE username = $1"
	row := r.db.QueryRow(query, username)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Phone, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := "SELECT id, first_name, last_name, username, phone, email, password, created_at, updated_at FROM users WHERE email = $1"
	row := r.db.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Phone, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	query := "SELECT id, first_name, last_name, username, phone, email, password, created_at, updated_at FROM users WHERE phone = $1"
	row := r.db.QueryRow(query, phone)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Phone, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) CreateUser(user *model.User) error {
	query := "INSERT INTO users (username, phone, password, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id"
	err := r.db.QueryRow(query, user.Username, user.Phone, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
