package grpc

import (
	"context"

	csatpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/delivery/proto"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/usecase"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type csatController struct {
	csatpb.UnimplementedCsatServiceServer
	csatUsecase usecase.ICsatUsecase
}

func NewCsatController(grpcServer *grpc.Server, csatUsecase usecase.ICsatUsecase) {
	controller := &csatController{
		csatUsecase: csatUsecase,
	}
	csatpb.RegisterCsatServiceServer(grpcServer, controller)
}

func (c *csatController) GetQuestions(ctx context.Context, _ *emptypb.Empty) (*csatpb.GetQuestionsResponse, error) {
	questions, err := c.csatUsecase.GetQuestions(ctx)
	if err != nil {
		return nil, err
	}

	var pbQuestions []*csatpb.Question
	for _, q := range questions {
		pbQuestions = append(pbQuestions, &csatpb.Question{
			Id:   q.ID.String(),
			Text: q.QuestionText,
		})
	}

	return &csatpb.GetQuestionsResponse{Questions: pbQuestions}, nil
}

func (c *csatController) CreateAnswer(ctx context.Context, req *csatpb.CreateAnswerRequest) (*emptypb.Empty, error) {
	questionID, err := uuid.Parse(req.GetQuestionId())
	if err != nil {
		return nil, usecase.ErrInvalidInput
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, usecase.ErrInvalidInput
	}

	answer := &model.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Rating:     model.RatingScale(req.GetRating()),
	}

	if err := c.csatUsecase.CreateAnswer(ctx, answer); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *csatController) GetStatistics(ctx context.Context, _ *emptypb.Empty) (*csatpb.GetStatisticsResponse, error) {
	stats, err := c.csatUsecase.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	var pbStats []*csatpb.QuestionStatistic
	for _, s := range stats.Questions {
		pbStats = append(pbStats, &csatpb.QuestionStatistic{
			QuestionId:    s.QuestionID.String(),
			AverageRating: s.AverageRating,
			AnswerCount:   int32(s.TotalResponses),
		})
	}

	var totalAnswers int
	var weightedSum float64
	for _, q := range stats.Questions {
		totalAnswers += q.TotalResponses
		weightedSum += float64(q.TotalResponses) * q.AverageRating
	}

	var overallAverage float64
	if totalAnswers > 0 {
		overallAverage = weightedSum / float64(totalAnswers)
	}

	return &csatpb.GetStatisticsResponse{
		Statistics: &csatpb.FullStatistics{
			QuestionStatistics: pbStats,
			AverageRating:      overallAverage,
		},
	}, nil
}

func (c *csatController) GetUserActivity(ctx context.Context, req *csatpb.GetUserActivityRequest) (*csatpb.GetUserActivityResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, usecase.ErrInvalidInput
	}

	activity, err := c.csatUsecase.GetUserActivity(ctx, userID)
	if err != nil {
		return nil, err
	}

	avgRating, err := c.csatUsecase.GetUserAverageRating(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &csatpb.GetUserActivityResponse{
		Activity: &csatpb.UserActivity{
			UserId:        activity.UserID.String(),
			TotalAnswers:  int32(activity.ResponsesCount),
			AverageRating: avgRating,
		},
	}, nil
}
