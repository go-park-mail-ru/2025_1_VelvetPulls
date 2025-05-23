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

func (uc *ContactUsecase) SearchContacts(ctx context.Context, userIDStr string, query string) ([]model.Contact, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, model.ErrValidation
	}

	metrics.IncBusinessOp("search_contacts")
	return uc.repo.SearchContacts(ctx, userID, query)
}
