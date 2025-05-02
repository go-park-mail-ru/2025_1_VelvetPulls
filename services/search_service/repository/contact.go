package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/model"
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
			u.first_name, 
			u.last_name, 
			u.avatar_path 
		FROM contact c
		JOIN public.user u ON c.contact_id = u.id
		WHERE c.user_id = $1 
			AND (
				u.username ILIKE $2 OR
				u.first_name ILIKE $2 OR 
				u.last_name ILIKE $2
			)`

	rows, err := r.db.QueryContext(ctx, querySQL, userID, "%"+query+"%")
	if err != nil {
		return nil, ErrSearchContacts
	}
	defer rows.Close()

	var contacts []model.Contact
	for rows.Next() {
		var c model.Contact
		var firstName, lastName, avatarPath sql.NullString

		err := rows.Scan(
			&c.ID,
			&c.Username,
			&firstName,
			&lastName,
			&avatarPath,
		)
		if err != nil {
			return nil, err
		}

		if firstName.Valid {
			c.FirstName = &firstName.String
		}
		if lastName.Valid {
			c.LastName = &lastName.String
		}
		if avatarPath.Valid {
			c.AvatarURL = &avatarPath.String
		}

		contacts = append(contacts, c)
	}
	return contacts, nil
}
