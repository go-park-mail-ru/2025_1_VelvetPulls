package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/server"
	_ "github.com/lib/pq"
)

func main() {
	config.Init()

	// Инициализация подключения к PostgreSQL
	dbConn, err := sql.Open("postgres", config.GetPostgresDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	dbConn.SetMaxOpenConns(11)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(30 * time.Minute)

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	defer dbConn.Close()

	// Создание и запуск сервера
	s := server.NewServer(dbConn)
	if err := s.Run(":8083"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
