package main

import (
	"database/sql"
	"log"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/server"
	_ "github.com/lib/pq"
)

func main() {
	config.Init()

	dbConn, err := sql.Open("postgres", config.GetPostgresDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	defer dbConn.Close()

	log.Printf("Starting server on :8082")
	s := server.NewServer(dbConn)
	if err := s.Run(":8082"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
