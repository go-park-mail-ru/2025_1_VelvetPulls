package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/repository"
)

type UserUsecase struct {
	repo repository.UserRepo
}

func NewUserUsecase(repo repository.UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) SearchUsers(ctx context.Context, query string) ([]model.UserProfile, error) {
	if len(query) < 3 {
		return nil, model.ErrValidation
	}

	metrics.IncBusinessOp("search_users")
	return uc.repo.SearchUsers(ctx, query)
}
