package usecase

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/google/uuid"
)

type IContactUsecase interface {
	GetUserContacts(ctx context.Context, userID uuid.UUID) (*[]model.Contact, error)
	AddUserContact(ctx context.Context, userID, contactID uuid.UUID) error
	RemoveUserContact(ctx context.Context, userID, contactID uuid.UUID) error
}
type ContactUsecase struct {
	contactRepo repository.IContactRepo
}

func NewContactUsecase(contactRepo repository.IContactRepo) IContactUsecase {
	return &ContactUsecase{contactRepo: contactRepo}
}

func (uc *ContactUsecase) GetUserContacts(ctx context.Context, userID uuid.UUID) (*[]model.Contact, error) {
	return uc.contactRepo.GetContacts(ctx, userID)
}

func (uc *ContactUsecase) AddUserContact(ctx context.Context, userID, contactID uuid.UUID) error {
	if userID == contactID {
		return errors.New("cannot add yourself as a contact")
	}
	return uc.contactRepo.AddContact(ctx, userID, contactID)
}

func (uc *ContactUsecase) RemoveUserContact(ctx context.Context, userID, contactID uuid.UUID) error {
	return uc.contactRepo.DeleteContact(ctx, userID, contactID)
}
