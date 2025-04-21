package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	_ "github.com/go-park-mail-ru/2025_1_VelvetPulls/docs"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/server"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
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

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	defer dbConn.Close()

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.Redis.Password,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis:", err)
	}

	defer redisClient.Close()

	log.Printf("Starting server on %s", config.PORT)

	s := server.NewServer(dbConn, redisClient, nc)
	if err := s.Run(config.PORT); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
