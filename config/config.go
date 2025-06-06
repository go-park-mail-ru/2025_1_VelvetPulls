package config

import (
	"fmt"
	"os"
	"time"
)

var (
	PORT           = ":8080"
	UPLOAD_DIR     = "./uploads/"
	LOG_DIR        = "./logs/"
	MAX_FILE_SIZE  = int64(2 << 20) // 2 MB (2,097,152 байт)
	CookieDuration = 3 * time.Hour
)

var CSRF = struct {
	CsrfAuthKey  string
	IsProduction bool
}{
	CsrfAuthKey:  "32-byte-long-auth-key-here",
	IsProduction: false,
}

var Cors = struct {
	AllowedOrigin  string
	AllowedMethods string
	AllowedHeaders string
}{
	AllowedOrigin:  "http://localhost:8088",
	AllowedMethods: "GET, POST, PUT, DELETE",
	AllowedHeaders: "Content-Type, Authorization, X-CSRF-Token, Access-Control-Allow-Credentials, enctype",
}

var (
	NATSURL  = "nats://nats:4222"
	NATSUser = ""
	NATSPass = ""
)

var Postgre = struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}{}

var Minio = struct {
	Host      string
	Port      string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}{}

var Redis = struct {
	Host     string
	Port     string
	Password string
}{}

func Init() {
	Postgre.User = os.Getenv("DATABASE_USER")
	Postgre.Password = os.Getenv("DATABASE_PASS")
	Postgre.Host = os.Getenv("DATABASE_HOST")
	Postgre.Port = os.Getenv("DATABASE_PORT")
	Postgre.DBName = os.Getenv("DATABASE_NAME")
	Postgre.SSLMode = os.Getenv("DATABASE_SSLMODE")

	Redis.Host = os.Getenv("REDIS_HOST")
	Redis.Port = os.Getenv("REDIS_PORT")
	Redis.Password = os.Getenv("REDIS_PASSWORD")

	Minio.Host = os.Getenv("MINIO_HOST")
	Minio.Port = os.Getenv("MINIO_PORT")
	Minio.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	Minio.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	Minio.UseSSL = os.Getenv("MINIO_USE_SSL") == "true"
	Minio.Bucket = os.Getenv("MINIO_BUCKET")
}

func GetPostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		Postgre.User, Postgre.Password, Postgre.Host, Postgre.Port, Postgre.DBName, Postgre.SSLMode)
}

func GetMinioEndpoint() string {
	return fmt.Sprintf("%s:%s", Minio.Host, Minio.Port)
}

func GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", Redis.Host, Redis.Port)
}
