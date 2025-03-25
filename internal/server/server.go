package server

import (
	"database/sql"
	"net/http"
	"time"

	delivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

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

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	// документация Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	// Repository
	sessionRepo := repository.NewSessionRepo(s.redisClient)
	authRepo := repository.NewauthRepo(s.dbConn)
	chatRepo := repository.NewChatRepo(s.dbConn)

	// Usecase
	authUsecase := usecase.NewAuthUsecase(authRepo, sessionRepo)
	chatUsecase := usecase.NewChatUsecase(sessionRepo, chatRepo)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo)

	// Controller
	delivery.NewAuthController(r, authUsecase)
	delivery.NewChatController(r, chatUsecase, sessionUsecase)

	httpServer := &http.Server{
		Handler:      middleware.CorsMiddleware(r),
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return httpServer.ListenAndServe()
}
