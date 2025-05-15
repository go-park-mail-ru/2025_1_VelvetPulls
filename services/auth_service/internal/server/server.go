package server

import (
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	grpcDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/delivery/grpc"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/usecase"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
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
	Run(address string) error
}

type Server struct {
	dbConn      *sql.DB
	redisClient *redis.Client
}

func NewServer(dbConn *sql.DB, redisClient *redis.Client) IServer {
	return &Server{dbConn: dbConn, redisClient: redisClient}
}

func (s *Server) Run(address string) error {
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

	// Repos
	sessionRepo := repository.NewSessionRepo(s.redisClient)
	authRepo := repository.NewAuthRepo(s.dbConn)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(authRepo, sessionRepo)
	sessionUsecase := usecase.NewSessionUsecase(authRepo, sessionRepo)

	// gRPC server с метриками
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			middleware.AccessLogInterceptor(),
			metrics.GrpcErrorCounterInterceptor(),
			metrics.GrpcHitCounterInterceptor(),
			metrics.GrpcTimingInterceptor(),
		),
	)

	grpcDelivery.NewAuthController(grpcServer, authUsecase, sessionUsecase)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	utils.Logger.Info("gRPC server started on " + address)
	return grpcServer.Serve(lis)
}
