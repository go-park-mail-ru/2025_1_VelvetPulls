package server

import (
	"net/http"
	"time"

	handler "github.com/go-park-mail-ru/2025_1_VelvetPulls/handler"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/middleware"
	repository "github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
	service "github.com/go-park-mail-ru/2025_1_VelvetPulls/service"
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

	// Подготовка Repository
	userRepo := repository.NewUserRepository()

	// Подготовка Service
	userService := service.NewUserService(userRepo, nil) // TODO: заменить nil на репозиторий сессии

	// Подготовка Handler
	userHandler := handler.NewUserHandler(r, userService)

	// Ручки роутера
	r.HandleFunc("/register/", userHandler.Register).Methods(http.MethodPost)

	// документация Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	// CORS
	handlerWithCORS := middleware.CorsMiddleware(r)

	httpServer := &http.Server{
		Handler:      handlerWithCORS,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return httpServer.ListenAndServe()
}
