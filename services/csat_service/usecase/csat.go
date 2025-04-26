package usecase

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ICsatUsecase interface {
	GetQuestions(ctx context.Context) ([]*model.Question, error)
	CreateAnswer(ctx context.Context, answer *model.Answer) error
	GetStatistics(ctx context.Context) (*model.FullStatistics, error)
	GetUserActivity(ctx context.Context, userID uuid.UUID) (*model.UserActivity, error)
	GetUserAverageRating(ctx context.Context, userID uuid.UUID) (float64, error)
}

type csatUsecase struct {
	repo repository.ICsatRepository
}

func NewCsatUsecase(repo repository.ICsatRepository) ICsatUsecase {
	return &csatUsecase{
		repo: repo,
	}
}

func (u *csatUsecase) GetQuestions(ctx context.Context) ([]*model.Question, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	questions, err := u.repo.GetQuestions(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Warn("No active questions found")
			return nil, ErrNotFound
		}
		logger.Error("Failed to get questions", zap.Error(err))
		return nil, ErrInternalServerError
	}

	return questions, nil
}

func (u *csatUsecase) CreateAnswer(ctx context.Context, answer *model.Answer) error {
	logger := utils.GetLoggerFromCtx(ctx)

	if answer.QuestionID == uuid.Nil || answer.UserID == uuid.Nil || answer.Rating < 1 || answer.Rating > 5 {
		logger.Warn("Invalid answer data",
			zap.Any("question_id", answer.QuestionID),
			zap.Any("user_id", answer.UserID),
			zap.Any("rating", answer.Rating))
		return ErrInvalidInput
	}

	err := u.repo.CreateAnswer(ctx, answer)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidInput):
			logger.Warn("Invalid input for answer creation", zap.Error(err))
			return ErrInvalidInput
		case errors.Is(err, repository.ErrDatabaseOperation):
			logger.Error("Failed to create answer", zap.Error(err))
			return ErrInternalServerError
		default:
			logger.Error("Unexpected error creating answer", zap.Error(err))
			return ErrInternalServerError
		}
	}

	return nil
}

func (u *csatUsecase) GetStatistics(ctx context.Context) (*model.FullStatistics, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	stats, err := u.repo.GetStatistics(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Warn("No statistics found")
			return nil, ErrNotFound
		}
		logger.Error("Failed to get statistics", zap.Error(err))
		return nil, ErrInternalServerError
	}

	return stats, nil
}

func (u *csatUsecase) GetUserActivity(ctx context.Context, userID uuid.UUID) (*model.UserActivity, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	if userID == uuid.Nil {
		logger.Warn("Invalid user ID")
		return nil, ErrInvalidInput
	}

	activity, err := u.repo.GetUserActivity(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Warn("User activity not found", zap.String("user_id", userID.String()))
			return nil, ErrNotFound
		}
		logger.Error("Failed to get user activity",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, ErrInternalServerError
	}

	return activity, nil
}

func (u *csatUsecase) GetUserAverageRating(ctx context.Context, userID uuid.UUID) (float64, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	if userID == uuid.Nil {
		logger.Warn("Invalid user ID")
		return 0, ErrInvalidInput
	}

	avgRating, err := u.repo.GetUserAverageRating(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user average rating",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return 0, ErrInternalServerError
	}

	return avgRating, nil
}
