package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/repository"
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

	return uc.repo.SearchUsers(ctx, query)
}
