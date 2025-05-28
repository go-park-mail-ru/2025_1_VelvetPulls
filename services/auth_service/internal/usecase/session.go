package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/repository"
	"github.com/google/uuid"
)

type ISessionUsecase interface {
	CheckLogin(ctx context.Context, token string) (*model.User, error)
}

type SessionUsecase struct {
	sessionRepo repository.ISessionRepo
	authRepo    repository.IAuthRepo
}

func NewSessionUsecase(authRepo repository.IAuthRepo, sessionRepo repository.ISessionRepo) ISessionUsecase {
	return &SessionUsecase{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *SessionUsecase) CheckLogin(ctx context.Context, token string) (*model.User, error) {
	userIDStr, err := uc.sessionRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err // возможно, стоит обернуть в свою ошибку: ErrInvalidUUID
	}

	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	metrics.IncBusinessOp("check_login")
	return user, nil
}
