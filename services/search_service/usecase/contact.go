package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/repository"
	"github.com/google/uuid"
)

type ContactUsecase struct {
	repo repository.ContactRepo
}

func NewContactUsecase(repo repository.ContactRepo) *ContactUsecase {
	return &ContactUsecase{repo: repo}
}

func (uc *ContactUsecase) SearchContacts(ctx context.Context, userIDStr string, query string) ([]model.Contact, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, model.ErrValidation
	}

	return uc.repo.SearchContacts(ctx, userID, query)
}
