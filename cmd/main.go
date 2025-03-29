package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-park-mail-ru/2025_1_VelvetPulls/docs"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// @title Keftegram backend API
// @version 1.0
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к БД
	postgreAddr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASS"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	dbConn, err := sql.Open("postgres", postgreAddr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	defer dbConn.Close()

	// Подключение к Redis
	redisAddr := fmt.Sprintf("localhost:%s", os.Getenv("REDIS_PORT"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis:", err)
	}

	defer redisClient.Close()

	// можно ввести свой порт при запуске
	addr := flag.String("addr", ":8080", "address for http server")

	log.Printf("Starting server on %s", *addr)

	s := server.NewServer(dbConn, redisClient)
	if err := s.Run(*addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
