package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/server"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.Redis.Password,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis:", err)
	}

	defer redisClient.Close()

	log.Printf("Starting server on %s", config.PORT)
	s := server.NewServer(dbConn, redisClient)
	if err := s.Run(":8081"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
