package server

import (
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	grpcDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/delivery/grpc"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/usecase"
	search "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func setupLogger() (*os.File, error) {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	utils.InitLoggerWithFile(logFile)
	utils.Logger.Info("Logger initialized")
	return logFile, nil
}

type IServer interface {
	Run(port string) error
}

type Server struct {
	dbConn *sql.DB
}

func NewServer(dbConn *sql.DB) IServer {
	return &Server{dbConn: dbConn}
}

func (s *Server) Run(port string) error {
	logFile, err := setupLogger()
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Запускаем HTTP сервер для метрик
	go func() {
		metricsRouter := http.NewServeMux()
		metricsRouter.Handle("/metrics", promhttp.Handler())

		metricsServer := &http.Server{
			Addr:    ":9091",
			Handler: metricsRouter,
		}

		utils.Logger.Info("Starting metrics server on :9091")
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Error("Metrics server error",
				zap.Error(err),
			)
		}
	}()

	// Инициализация репозиториев
	chatRepo := repository.NewChatRepo(s.dbConn)
	contactRepo := repository.NewContactRepo(s.dbConn)
	userRepo := repository.NewUserRepo(s.dbConn)
	messageRepo := repository.NewMessageRepo(s.dbConn)

	// Инициализация usecases
	chatUC := usecase.NewChatUsecase(*chatRepo)
	contactUC := usecase.NewContactUsecase(*contactRepo)
	userUC := usecase.NewUserUsecase(*userRepo)
	messageUC := usecase.NewMessageUsecase(*messageRepo)

	// Настройка gRPC сервера
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			middleware.AccessLogInterceptor(),
			metrics.GrpcErrorCounterInterceptor(),
			metrics.GrpcHitCounterInterceptor(),
			metrics.GrpcTimingInterceptor(),
		),
	)

	// Регистрация обработчиков
	handler := grpcDelivery.NewChatHandler(*chatUC, *contactUC, *userUC, *messageUC)
	search.RegisterChatServiceServer(grpcServer, handler)

	// Запуск сервера
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	utils.Logger.Info("gRPC chat server started on port " + port)
	return grpcServer.Serve(lis)
}
