package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/repository"
	"github.com/google/uuid"
)

type ContactUsecase struct {
	repo repository.ContactRepo
}

func NewContactUsecase(repo repository.ContactRepo) *ContactUsecase {
	return &ContactUsecase{repo: repo}
}

func (uc *ContactUsecase) SearchContacts(ctx context.Context, userIDStr string, query string) ([]model.UserProfile, []model.Contact, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, nil, model.ErrValidation
	}

	metrics.IncBusinessOp("search_contacts")
	users, err := uc.repo.SearchUsers(ctx, query)
	if err != nil {
		return nil, nil, err
	}

	contacts, err := uc.repo.SearchContacts(ctx, userID, query)
	if err != nil {
		return nil, nil, err
	}
	return users, contacts, nil
}
