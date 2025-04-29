package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) SearchUsers(ctx context.Context, query string) ([]model.UserProfile, error) {
	querySQL := `
		SELECT 
			id, 
			username, 
			first_name, 
			last_name, 
			avatar_path 
		FROM public.user
		WHERE username ILIKE $1 
			OR first_name ILIKE $1 
			OR last_name ILIKE $1
		LIMIT 50
	`

	rows, err := r.db.QueryContext(ctx, querySQL, "%"+query+"%")
	if err != nil {
		return nil, ErrSearchUsers
	}
	defer rows.Close()

	var users []model.UserProfile
	for rows.Next() {
		var u model.UserProfile
		var firstName, lastName, avatarPath sql.NullString

		err := rows.Scan(
			&u.Username,
			&firstName,
			&lastName,
			&avatarPath,
		)
		if err != nil {
			return nil, err
		}

		if firstName.Valid {
			u.FirstName = &firstName.String
		}
		if lastName.Valid {
			u.LastName = &lastName.String
		}
		if avatarPath.Valid {
			u.AvatarPath = &avatarPath.String
		}

		users = append(users, u)
	}
	return users, nil
}
