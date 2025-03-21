package server

import (
	"net/http"
	"os"
	"time"

	delivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/utils"
	"github.com/gorilla/mux"
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

// TODO: добавить объекты для подключения к бд
type Server struct {
	dbConn *int
}

func NewServer(dbConn *int) *Server {
	return &Server{dbConn: dbConn}
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
	sessionRepo := repository.NewSessionRepo()
	userRepo := repository.NewUserRepo()
	chatRepo := repository.NewChatRepo()

	// Usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo)
	chatUsecase := usecase.NewChatUsecase(sessionRepo, chatRepo)

	// Controller
	delivery.NewAuthController(r, authUsecase)
	delivery.NewChatController(r, chatUsecase)

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
