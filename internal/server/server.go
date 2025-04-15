package server

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	httpDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	websocketDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/websocket"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
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

	// ===== Root Router =====
	mainRouter := mux.NewRouter()

	// ===== API Subrouter =====
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()

	// Repository
	sessionRepo := repository.NewSessionRepo(s.redisClient)
	userRepo := repository.NewUserRepo(s.dbConn)
	chatRepo := repository.NewChatRepo(s.dbConn)
	contactRepo := repository.NewContactRepo(s.dbConn)
	messageRepo := repository.NewMessageRepo(s.dbConn)

	// Usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo)
	websocketUsecase := usecase.NewWebsocketUsecase(chatRepo)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, chatRepo, websocketUsecase)
	chatUsecase := usecase.NewChatUsecase(chatRepo)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	contactUsecase := usecase.NewContactUsecase(contactRepo)

	// Controllers
	httpDelivery.NewAuthController(apiRouter, authUsecase)
	httpDelivery.NewChatController(apiRouter, chatUsecase, sessionUsecase)
	httpDelivery.NewUserController(apiRouter, userUsecase, sessionUsecase)
	httpDelivery.NewMessageController(apiRouter, messageUsecase, sessionUsecase)
	httpDelivery.NewContactController(apiRouter, contactUsecase, sessionUsecase)

	// ===== WebSocket =====
	websocketDelivery.NewWebsocketController(mainRouter, sessionUsecase, websocketUsecase)

	// ===== Uploads =====
	uploadsRouter := mainRouter.PathPrefix("/uploads").Subrouter()
	httpDelivery.NewUploadsController(uploadsRouter)

	// Swagger
	apiRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	// CSRF
	apiRouter.HandleFunc("/csrf", httpDelivery.GetCSRFTokenHandler).Methods(http.MethodGet)

	// Middleware only for API
	apiRouter.Use(middleware.RequestIDMiddleware)
	apiRouter.Use(middleware.AccessLogMiddleware)

	handler := middleware.CorsMiddleware(mainRouter)
	// handler = middleware.CSRFMiddleware(config.CSRF.IsProduction, []byte(config.CSRF.CsrfAuthKey))(handler)

	// Server with CORS applied globally
	httpServer := &http.Server{
		Handler:      handler,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return httpServer.ListenAndServe()
}
