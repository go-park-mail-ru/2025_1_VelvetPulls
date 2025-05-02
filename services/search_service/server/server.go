package server

import (
	"database/sql"
	"net"
	"os"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	grpcDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/delivery/grpc"
	search "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/delivery/proto"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/usecase"
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
