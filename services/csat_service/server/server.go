package server

import (
	"database/sql"
	"net"
	"os"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	grpcDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/delivery/grpc"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/usecase"
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
	dbConn *sql.DB
}

func NewServer(dbConn *sql.DB) IServer {
	return &Server{dbConn: dbConn}
}

func (s *Server) Run(address string) error {
	logFile, err := setupLogger()
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Repos
	csatRepo := repository.NewCsatRepository(s.dbConn)
	// Usecases
	csatUsecase := usecase.NewCsatUsecase(csatRepo)
	// gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDInterceptor(),
			middleware.AccessLogInterceptor(),
		),
	)

	grpcDelivery.NewCsatController(grpcServer, csatUsecase)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return grpcServer.Serve(lis)
}
