package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type IContactUsecase interface {
	GetUserContacts(ctx context.Context, userID uuid.UUID) ([]model.Contact, error)
	AddUserContact(ctx context.Context, userID uuid.UUID, contactUsername string) (*model.Contact, error)
	RemoveUserContact(ctx context.Context, userID uuid.UUID, contactUsername string) error
}
type ContactUsecase struct {
	contactRepo repository.IContactRepo
}

func NewContactUsecase(contactRepo repository.IContactRepo) IContactUsecase {
	return &ContactUsecase{contactRepo: contactRepo}
}

func (uc *ContactUsecase) GetUserContacts(ctx context.Context, userID uuid.UUID) ([]model.Contact, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Fetching contacts")
	metrics.IncBusinessOp("get_contacts")
	return uc.contactRepo.GetContacts(ctx, userID)
}

func (uc *ContactUsecase) AddUserContact(ctx context.Context, userID uuid.UUID, contactUsername string) (*model.Contact, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Adding contact")
	metrics.IncBusinessOp("add_contact")

	contact, err := uc.contactRepo.AddContactByUsername(ctx, userID, contactUsername)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (uc *ContactUsecase) RemoveUserContact(ctx context.Context, userID uuid.UUID, contactUsername string) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Removing contact")
	metrics.IncBusinessOp("remove_contact")
	return uc.contactRepo.DeleteContactByUsername(ctx, userID, contactUsername)
}
