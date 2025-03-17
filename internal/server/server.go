package server

import (
	"net/http"
	"time"

	a "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http/auth"
	c "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http/chat"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase/chat"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// TODO: добавить объекты для подключения к бд
type Server struct {
	dbConn *int
}

func NewServer(dbConn *int) *Server {
	return &Server{dbConn: dbConn}
}

// TODO: подключиться к бд
func (s *Server) Run(address string) error {

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	// документация Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	// Repository
	sessionRepo := repository.NewSessionRepo()
	userRepo := repository.NewUserRepo()
	chatRepo := repository.NewChatRepo()

	// Usecase
	authUsecase := auth.NewAuthUsecase(userRepo, sessionRepo)
	chatUsecase := chat.NewChatUsecase(sessionRepo, chatRepo)

	// Controller
	a.NewAuthController(r, authUsecase)
	c.NewChatController(r, chatUsecase)

	httpServer := &http.Server{
		Handler:      middleware.CorsMiddleware(r),
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return httpServer.ListenAndServe()
}
