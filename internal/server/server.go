package server

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	httpDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	generatedAuth "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/delivery/proto"
	searchpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/delivery/proto"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	httpSwagger "github.com/swaggo/http-swagger"
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
	dbConn     *sql.DB
	authConn   *grpc.ClientConn
	searchConn *grpc.ClientConn
	nc         *nats.Conn
}

func NewServer(dbConn *sql.DB, authConn *grpc.ClientConn, searchConn *grpc.ClientConn, nc *nats.Conn) IServer {
	return &Server{dbConn: dbConn, authConn: authConn, searchConn: searchConn, nc: nc}
}

func (s *Server) Run(address string) error {
	logFile, err := setupLogger()
	if err != nil {
		return err
	}
	defer logFile.Close()

	// ===== Microservice usecase =====
	authClient := generatedAuth.NewAuthServiceClient(s.authConn)
	sessionClient := generatedAuth.NewSessionServiceClient(s.authConn)
	searchClient := searchpb.NewChatServiceClient(s.searchConn)
	// ===== Root Router =====
	mainRouter := mux.NewRouter()

	mainRouter.Handle("/metrics", promhttp.Handler())

	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.AccessLogMiddleware)
	mainRouter.Use(metrics.TimingHistogramMiddleware) // Замер времени выполнения
	mainRouter.Use(metrics.HitCounterMiddleware)      // Счетчик запросов
	mainRouter.Use(metrics.ErrorCounterMiddleware)    // Счетчик ошибок

	// ===== API Subrouter =====
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()

	// Repository
	userRepo := repository.NewUserRepo(s.dbConn)
	chatRepo := repository.NewChatRepo(s.dbConn)
	contactRepo := repository.NewContactRepo(s.dbConn)
	messageRepo := repository.NewMessageRepo(s.dbConn)

	// Usecase
	messageUsecase := usecase.NewMessageUsecase(messageRepo, chatRepo, s.nc)
	chatUsecase := usecase.NewChatUsecase(chatRepo, s.nc)
	userUsecase := usecase.NewUserUsecase(userRepo)
	contactUsecase := usecase.NewContactUsecase(contactRepo)

	// Controllers
	httpDelivery.NewAuthController(apiRouter, authClient, sessionClient)
	httpDelivery.NewChatController(apiRouter, chatUsecase, sessionClient)
	httpDelivery.NewUserController(apiRouter, userUsecase, sessionClient)
	httpDelivery.NewMessageController(apiRouter, messageUsecase, sessionClient)
	httpDelivery.NewContactController(apiRouter, contactUsecase, sessionClient)
	httpDelivery.NewSearchController(apiRouter, searchClient, sessionClient)

	// ===== Uploads =====
	uploadsRouter := mainRouter.PathPrefix("/uploads").Subrouter()
	httpDelivery.NewUploadsController(uploadsRouter)

	// Swagger
	apiRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	handler := middleware.CorsMiddleware(mainRouter)

	// Server with CORS applied globally
	httpServer := &http.Server{
		Handler:      handler,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return httpServer.ListenAndServe()
}
