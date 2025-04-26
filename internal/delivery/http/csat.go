package http

import (
	"context"
	"net/http"
	"time"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	csatpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/delivery/proto"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type csatController struct {
	csatClient csatpb.CsatServiceClient
}

func NewCsatController(r *mux.Router, csatClient csatpb.CsatServiceClient) {
	controller := &csatController{
		csatClient: csatClient,
	}

	r.Handle("/csat/questions", http.HandlerFunc(controller.GetQuestions)).Methods(http.MethodGet)
	r.Handle("/csat/answers", http.HandlerFunc(controller.CreateAnswer)).Methods(http.MethodPost)
	r.Handle("/csat/statistics", http.HandlerFunc(controller.GetStatistics)).Methods(http.MethodGet)
	r.Handle("/csat/activity", http.HandlerFunc(controller.GetUserActivity)).Methods(http.MethodGet)
}

func (c *csatController) GetQuestions(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := c.csatClient.GetQuestions(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error("gRPC GetQuestions error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, resp.Questions, true)
}

func (c *csatController) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	var answer model.CreateCsatAnswerRequest
	if err := utils.ParseJSONRequest(r, &answer); err != nil {
		logger.Warn("Invalid request data", zap.Error(err))
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request data", false)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := c.csatClient.CreateAnswer(ctx, &csatpb.CreateAnswerRequest{
		QuestionId: answer.QuestionID,
		Username:   answer.Username,
		Rating:     int32(answer.Rating),
	})
	if err != nil {
		logger.Error("gRPC CreateAnswer error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, "Answer created successfully", true)
}

func (c *csatController) GetStatistics(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := c.csatClient.GetStatistics(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Error("gRPC GetStatistics error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, resp.Statistics, true)
}

func (c *csatController) GetUserActivity(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		logger.Warn("Missing user_id in query")
		utils.SendJSONResponse(w, http.StatusBadRequest, "Missing user_id parameter", false)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := c.csatClient.GetUserActivity(ctx, &csatpb.GetUserActivityRequest{
		Username: userID,
	})
	if err != nil {
		logger.Error("gRPC GetUserActivity error", zap.Error(err))
		code, msg := apperrors.UnpackGrpcError(err)
		utils.SendJSONResponse(w, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, resp.Activity, true)
}
