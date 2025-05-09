package server

import (
	"net/http"
	"os"
	"time"

	middleware "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	generatedAuth "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	websocketDelivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/websocket_service/internal/delivery/websocket"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/websocket_service/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	nc       *nats.Conn
	authConn *grpc.ClientConn
}

func NewServer(nc *nats.Conn, authConn *grpc.ClientConn) IServer {
	return &Server{nc: nc, authConn: authConn}
}

func (s *Server) Run(address string) error {
	logFile, err := setupLogger()
	if err != nil {
		return err
	}
	defer logFile.Close()

	// ===== Microservice usecase =====
	sessionClient := generatedAuth.NewSessionServiceClient(s.authConn)

	mainRouter := mux.NewRouter()

	mainRouter.Handle("/metrics", promhttp.Handler())

	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.AccessLogMiddleware)

	// Usecases
	websocketUsecase := usecase.NewWebsocketUsecase(s.nc)

	// ===== WebSocket =====
	websocketDelivery.NewWebsocketController(mainRouter, sessionClient, websocketUsecase, s.nc)

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
