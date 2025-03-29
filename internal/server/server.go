package server

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	delivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/utils"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
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

// TODO: добавить объекты для подключения к бд
type Server struct {
	dbConn      *sql.DB
	redisClient *redis.Client
}

func NewServer(dbConn *sql.DB, redisClient *redis.Client) IServer {
	return &Server{dbConn: dbConn, redisClient: redisClient}
}

// TODO: подключиться к бд
func (s *Server) Run(address string) error {
	logFile, err := setupLogger()
	if err != nil {
		return err
	}
	defer logFile.Close()

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	// документация Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	// Repository
	sessionRepo := repository.NewSessionRepo(s.redisClient)
	userRepo := repository.NewUserRepo(s.dbConn)
	// chatRepo := repository.NewChatRepo(s.dbConn)

	// Usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo)
	// chatUsecase := usecase.NewChatUsecase(sessionRepo, chatRepo)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Controller
	delivery.NewAuthController(r, authUsecase)
	// delivery.NewChatController(r, chatUsecase, sessionUsecase)
	delivery.NewUserController(r, userUsecase, sessionUsecase)
	delivery.NewUploadsController(r)

	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.AccessLogMiddleware)

	httpServer := &http.Server{
		Handler:      r,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return httpServer.ListenAndServe()
}
