package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	PORT           = ":8080"
	UPLOAD_DIR     = "./uploads/"
	MAX_FILE_SIZE  = 2 * 1024 * 1024
	CookieDuration = 3 * time.Hour
)

var Cors = struct {
	AllowedOrigin  string
	AllowedMethods string
	AllowedHeaders string
}{
	AllowedOrigin:  "http://localhost:8081",
	AllowedMethods: "GET, POST, PUT, DELETE",
	AllowedHeaders: "Content-Type, Authorization",
}

var Postgre = struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}{}

var Redis = struct {
	Host     string
	Port     string
	Password string
}{}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found or error loading .env - using system environment variables")
	}

	Postgre.User = os.Getenv("DATABASE_USER")
	Postgre.Password = os.Getenv("DATABASE_PASS")
	Postgre.Host = os.Getenv("DATABASE_HOST")
	Postgre.Port = os.Getenv("DATABASE_PORT")
	Postgre.DBName = os.Getenv("DATABASE_NAME")
	Postgre.SSLMode = os.Getenv("DATABASE_SSLMODE")

	Redis.Host = os.Getenv("REDIS_HOST")
	Redis.Port = os.Getenv("REDIS_PORT")
	Redis.Password = os.Getenv("REDIS_PASSWORD")
}

func GetPostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		Postgre.User, Postgre.Password, Postgre.Host, Postgre.Port, Postgre.DBName, Postgre.SSLMode)
}

func GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", Redis.Host, Redis.Port)
}
