package main

import (
	"log"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/websocket_service/internal/server"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config.Init()

	nc, err := nats.Connect(
		config.NATSURL,
		nats.UserInfo(config.NATSUser, config.NATSPass),
		nats.Timeout(5*time.Second),       // таймаут подключения
		nats.ReconnectWait(2*time.Second), // интервал переподключения
		nats.MaxReconnects(10),            // максимальное число попыток переподключения
	)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	authConn, errAuth := grpc.NewClient("auth:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if errAuth != nil {
		log.Fatalf("Failed to connect to AuthService: %v", errAuth)
	}
	defer authConn.Close()

	log.Printf("Starting server on :8082")
	s := server.NewServer(nc, authConn)
	if err := s.Run(":8082"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
