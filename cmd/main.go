package main

import (
	"database/sql"
	"log"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	_ "github.com/go-park-mail-ru/2025_1_VelvetPulls/docs"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/server"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title Keftegram backend API
// @version 1.0
func main() {
	config.Init()

	// Подключение к БД
	dbConn, err := sql.Open("postgres", config.GetPostgresDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	defer dbConn.Close()

	authConn, errAuth := grpc.NewClient("auth:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	csatConn, errCsat := grpc.NewClient("csat:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if errAuth != nil {
		log.Fatalf("Failed to connect to AuthService: %v", errAuth)
	}
	defer authConn.Close()

	if errCsat != nil {
		log.Fatalf("Failed to connect to CsatService: %v", errCsat)
	}
	defer csatConn.Close()

	log.Printf("Starting server on %s", config.PORT)
	s := server.NewServer(dbConn, authConn, csatConn)
	if err := s.Run(config.PORT); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
