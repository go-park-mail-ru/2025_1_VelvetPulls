package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/google/uuid"
)

type ContactRepo struct {
	db *sql.DB
}

func NewContactRepo(db *sql.DB) *ContactRepo {
	return &ContactRepo{db: db}
}

func (r *ContactRepo) SearchContacts(ctx context.Context, userID uuid.UUID, query string) ([]model.Contact, error) {
	querySQL := `
		SELECT 
			u.id, 
			u.username, 
			u.name, 
			u.avatar_path 
		FROM contact c
		JOIN public.user u ON c.contact_id = u.id
		WHERE c.user_id = $1 
			AND (
				u.username ILIKE $2 OR
				u.name ILIKE $2 
			)`

	rows, err := r.db.QueryContext(ctx, querySQL, userID, "%"+query+"%")
	if err != nil {
		return nil, ErrSearchContacts
	}
	defer rows.Close()

	var contacts []model.Contact
	for rows.Next() {
		var c model.Contact
		var name, avatarPath sql.NullString

		err := rows.Scan(
			&c.ID,
			&c.Username,
			&name,
			&avatarPath,
		)
		if err != nil {
			return nil, err
		}

		if name.Valid {
			c.Name = &name.String
		}
		if avatarPath.Valid {
			c.AvatarURL = &avatarPath.String
		}

		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *ContactRepo) SearchUsers(ctx context.Context, query string) ([]model.UserProfile, error) {
	querySQL := `
        SELECT 
			id,
            username, 
            name, 
            avatar_path,
			birth_date
        FROM public.user
        WHERE 
            username ILIKE $1 OR
            name ILIKE $1 
        ORDER BY 
            username ILIKE $1 DESC,
            name ILIKE $1 DESC
        LIMIT 50`

	rows, err := r.db.QueryContext(ctx, querySQL, "%"+query+"%")
	if err != nil {
		return nil, ErrSearchUsers
	}
	defer rows.Close()

	var users []model.UserProfile
	for rows.Next() {
		var u model.UserProfile
		var name, avatarPath sql.NullString
		var bday sql.NullTime

		err := rows.Scan(
			&u.ID,
			&u.Username,
			&name,
			&avatarPath,
			&bday,
		)
		if err != nil {
			return nil, err
		}

		if bday.Valid {
			u.BirthDate = &bday.Time
		} else {
			u.BirthDate = nil
		}

		if name.Valid {
			u.Name = &name.String
		}
		if avatarPath.Valid {
			u.AvatarPath = &avatarPath.String
		}

		users = append(users, u)
	}
	return users, nil
}
