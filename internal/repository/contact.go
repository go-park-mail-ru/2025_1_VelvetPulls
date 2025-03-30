package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/google/uuid"
)

type IContactRepo interface {
	GetContacts(ctx context.Context, userID uuid.UUID) (*[]model.Contact, error)
	AddContact(ctx context.Context, userID, contactID uuid.UUID) error
	DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error
}

type contactRepo struct {
	db *sql.DB
}

func NewContactRepo(db *sql.DB) IContactRepo {
	return &contactRepo{db: db}
}

func (r *contactRepo) GetContacts(ctx context.Context, userID uuid.UUID) (*[]model.Contact, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, u.username, u.avatar_path
		FROM public.contact c
		JOIN public.user u ON c.contact_id = u.id
		WHERE c.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []model.Contact
	for rows.Next() {
		var contact model.Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.Username, &contact.AvatarURL); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &contacts, nil
}

func (r *contactRepo) AddContact(ctx context.Context, userID, contactID uuid.UUID) error {
	if userID == contactID {
		return errors.New("cannot add yourself as a contact")
	}

	query := `INSERT INTO public.contact (user_id, contact_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, contactID)
	return err
}

func (r *contactRepo) DeleteContact(ctx context.Context, userID, contactID uuid.UUID) error {
	query := `DELETE FROM public.contact WHERE user_id = $1 AND contact_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, contactID)
	return err
}
