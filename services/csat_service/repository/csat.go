package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ICsatRepository interface {
	GetQuestions(ctx context.Context) ([]*model.Question, error)
	CreateAnswer(ctx context.Context, answer *model.Answer) error
	GetStatistics(ctx context.Context) (*model.FullStatistics, error)
	GetUserActivity(ctx context.Context, username string) (*model.UserActivity, error)
	GetUserAverageRating(ctx context.Context, username string) (float64, error)
}

type psqlCsatRepository struct {
	db *sql.DB
}

func NewCsatRepository(db *sql.DB) ICsatRepository {
	return &psqlCsatRepository{
		db: db,
	}
}

func (r *psqlCsatRepository) GetQuestions(ctx context.Context) ([]*model.Question, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	query := `
		SELECT id, title, question_text, is_active, created_at, updated_at
		FROM csat.question
		WHERE is_active = TRUE
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("GetQuestions query failed", zap.Error(err))
		return nil, ErrDatabaseOperation
	}
	defer rows.Close()

	var questions []*model.Question
	for rows.Next() {
		var q model.Question
		err := rows.Scan(
			&q.ID,
			&q.Title,
			&q.QuestionText,
			&q.IsActive,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			logger.Error("Scan question failed", zap.Error(err))
			return nil, ErrDatabaseOperation
		}
		questions = append(questions, &q)
	}

	if len(questions) == 0 {
		return nil, ErrNotFound
	}

	return questions, nil
}

func (r *psqlCsatRepository) CreateAnswer(ctx context.Context, answer *model.Answer) error {
	logger := utils.GetLoggerFromCtx(ctx)

	if answer.QuestionID == uuid.Nil || answer.Rating < 1 || answer.Rating > 5 || answer.Username == "" {
		logger.Warn("Invalid answer data")
		return ErrInvalidInput
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Begin transaction failed", zap.Error(err))
		return ErrDatabaseOperation
	}
	defer tx.Rollback()

	// 1. Сохраняем ответ
	answerQuery := `
		INSERT INTO csat.answer (question_id, username, rating, feedback)
		VALUES ($1, $2, $3::rating_scale, $4)
	`
	_, err = tx.ExecContext(ctx, answerQuery,
		answer.QuestionID,
		answer.Username,
		fmt.Sprintf("%d", answer.Rating),
		answer.Feedback,
	)
	if err != nil {
		logger.Error("CreateAnswer failed", zap.Error(err))
		return ErrDatabaseOperation
	}

	// 2. Обновляем активность пользователя
	activityQuery := `
		INSERT INTO csat.user_activity (username, last_response_at, responses_count)
		VALUES ($1, NOW(), 1)
		ON CONFLICT (username) DO UPDATE
		SET last_response_at = NOW(),
		    responses_count = user_activity.responses_count + 1
	`
	_, err = tx.ExecContext(ctx, activityQuery, answer.Username)
	if err != nil {
		logger.Error("UpdateUserActivity failed", zap.Error(err))
		return ErrDatabaseOperation
	}

	if err = tx.Commit(); err != nil {
		logger.Error("Commit transaction failed", zap.Error(err))
		return ErrDatabaseOperation
	}

	return nil
}

func (r *psqlCsatRepository) GetStatistics(ctx context.Context) (*model.FullStatistics, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	// Получаем все активные вопросы
	questions, err := r.GetQuestions(ctx)
	if err != nil {
		return nil, err
	}

	stats := &model.FullStatistics{
		Questions: make([]model.QuestionStatistics, 0, len(questions)),
	}

	for _, q := range questions {
		// Получаем распределение оценок
		distrQuery := `
			SELECT rating, COUNT(*) as count
			FROM csat.answer
			WHERE question_id = $1
			GROUP BY rating
			ORDER BY rating
		`

		rows, err := r.db.QueryContext(ctx, distrQuery, q.ID)
		if err != nil {
			logger.Error("Get rating distribution failed", zap.Error(err))
			return nil, ErrDatabaseOperation
		}

		var (
			total     int
			sum       int
			distrib   []model.RatingDistribution
			ratingStr string
			count     int
		)

		for rows.Next() {
			if err := rows.Scan(&ratingStr, &count); err != nil {
				rows.Close()
				return nil, ErrDatabaseOperation
			}

			rating := int(ratingStr[0] - '0')
			distrib = append(distrib, model.RatingDistribution{
				Rating: model.RatingScale(rating),
				Count:  count,
			})

			total += count
			sum += rating * count
		}
		rows.Close()

		// Получаем комментарии
		commentsQuery := `
			SELECT id, username, rating, feedback, created_at
			FROM csat.answer
			WHERE question_id = $1 AND feedback IS NOT NULL
			ORDER BY created_at DESC
		`

		commentRows, err := r.db.QueryContext(ctx, commentsQuery, q.ID)
		if err != nil {
			logger.Error("Get comments failed", zap.Error(err))
			return nil, ErrDatabaseOperation
		}

		var comments []*model.Answer
		for commentRows.Next() {
			var a model.Answer
			var ratingStr string
			if err := commentRows.Scan(
				&a.ID,
				&a.Username,
				&ratingStr,
				&a.Feedback,
				&a.CreatedAt,
			); err != nil {
				commentRows.Close()
				return nil, ErrDatabaseOperation
			}

			a.Rating = model.RatingScale(int(ratingStr[0] - '0'))
			a.QuestionID = q.ID
			comments = append(comments, &a)
		}
		commentRows.Close()

		var avg float64
		if total > 0 {
			avg = float64(sum) / float64(total)
		}

		stats.Questions = append(stats.Questions, model.QuestionStatistics{
			QuestionID:     q.ID,
			QuestionText:   q.QuestionText,
			AverageRating:  avg,
			TotalResponses: total,
			Distribution:   distrib,
			Comments:       comments,
		})
	}

	if len(stats.Questions) == 0 {
		return nil, ErrNotFound
	}

	return stats, nil
}

func (r *psqlCsatRepository) GetUserActivity(ctx context.Context, username string) (*model.UserActivity, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	query := `
		SELECT username, last_response_at, responses_count
		FROM csat.user_activity
		WHERE username = $1
	`

	var activity model.UserActivity
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&activity.Username,
		&activity.LastResponseAt,
		&activity.ResponsesCount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		logger.Error("GetUserActivity failed", zap.Error(err))
		return nil, ErrDatabaseOperation
	}

	return &activity, nil
}

func (r *psqlCsatRepository) GetUserAverageRating(ctx context.Context, username string) (float64, error) {
	logger := utils.GetLoggerFromCtx(ctx)

	query := `
        SELECT AVG(rating::integer)
        FROM csat.answer
        WHERE username = $1
    `

	var avgRating *float64
	err := r.db.QueryRowContext(ctx, query, username).Scan(&avgRating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		logger.Error("GetUserAverageRating failed", zap.Error(err))
		return 0, ErrDatabaseOperation
	}

	if avgRating == nil {
		return 0, nil
	}

	return *avgRating, nil
}
