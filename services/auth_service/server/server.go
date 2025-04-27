package server

import (
	"database/sql"
	"net"
	"os"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	grpcDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/delivery/grpc"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/usecase"
	"github.com/redis/go-redis/v9"
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

	// Repos
	sessionRepo := repository.NewSessionRepo(s.redisClient)
	authRepo := repository.NewAuthRepo(s.dbConn)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(authRepo, sessionRepo)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo)

	// gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			middleware.AccessLogInterceptor(),
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
