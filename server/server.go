package server

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/middleware"
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
	r := mux.NewRouter()

	// Подготовка Repository

	// Подготовка Service

	// Подготовка Handler

	// Ручки роутера

	// документация Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET")

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
