package server

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	httpDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	websocketDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/websocket"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	minioRepo "github.com/go-park-mail-ru/2025_1_VelvetPulls/minio"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	minioConfig *MinioConfig
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func NewServer(dbConn *sql.DB, redisClient *redis.Client, minioConfig *MinioConfig) IServer {
	return &Server{
		dbConn:      dbConn,
		redisClient: redisClient,
		minioConfig: minioConfig,
	}
}
func (s *Server) initMinio() (*minio.Client, error) {
	// Initialize Minio client
	minioClient, err := minio.New(s.minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.minioConfig.AccessKey, s.minioConfig.SecretKey, ""),
		Secure: s.minioConfig.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Check if bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, s.minioConfig.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, s.minioConfig.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return minioClient, nil
}
func (s *Server) Run(address string) error {
	logFile, err := setupLogger()
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Initialize Minio client
	minioClient, err := s.initMinio()
	if err != nil {
		utils.Logger.Error("Failed to initialize Minio: " + err.Error())
		return err
	}
	utils.Logger.Info("Minio client initialized successfully")

	// ===== Root Router =====
	mainRouter := mux.NewRouter()

	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.AccessLogMiddleware)

	// ===== API Subrouter =====
	apiRouter := mainRouter.PathPrefix("/api").Subrouter()

	// Repository
	sessionRepo := repository.NewSessionRepo(s.redisClient)
	userRepo := repository.NewUserRepo(s.dbConn)
	chatRepo := repository.NewChatRepo(s.dbConn)
	contactRepo := repository.NewContactRepo(s.dbConn)
	messageRepo := repository.NewMessageRepo(s.dbConn)
	fileRepo := minioRepo.NewFileRepository(minioClient, s.minioConfig.Bucket)

	// Usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo)
	websocketUsecase := usecase.NewWebsocketUsecase(chatRepo)
	messageUsecase := usecase.NewMessageUsecase(messageRepo, chatRepo, websocketUsecase)
	chatUsecase := usecase.NewChatUsecase(chatRepo, websocketUsecase)
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
	httpDelivery.NewUploadsController(uploadsRouter, fileRepo)

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
