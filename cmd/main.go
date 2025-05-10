package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	_ "github.com/go-park-mail-ru/2025_1_VelvetPulls/docs"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/server"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title Keftegram backend API
// @version 1.0
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

	// Подключение к БД
	dbConn, err := sql.Open("postgres", config.GetPostgresDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	dbConn.SetMaxOpenConns(17)
	dbConn.SetMaxIdleConns(8)
	dbConn.SetConnMaxLifetime(30 * time.Minute)

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	defer dbConn.Close()

	authConn, errAuth := grpc.NewClient("auth:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if errAuth != nil {
		log.Fatalf("Failed to connect to AuthService: %v", errAuth)
	}
	defer authConn.Close()

	log.Printf("Starting server on %s", config.PORT)
	searchConn, err := grpc.NewClient("search:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to ChatService: %v", err)
	}
	defer searchConn.Close()

	s := server.NewServer(dbConn, authConn, searchConn, nc)
	if err := s.Run(config.PORT); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
