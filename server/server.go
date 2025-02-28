package server

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/handler"
	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/middleware"
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

	// Ручки роутера
	r.HandleFunc("/register/", handler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login/", handler.Login).Methods(http.MethodPost)
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
