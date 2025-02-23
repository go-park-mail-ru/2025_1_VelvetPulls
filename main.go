package main

import (
	"flag"
	"log"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/server"
)

func main() {
	// TODO: сделать подтягивание параметорв из env

	// TODO: подключение бд
	dbConn := new(int)

	// можно ввести свой порт при запуске
	addr := flag.String("addr", ":8080", "address for http server")
	log.Printf("Starting server on %s", *addr)
	s := server.NewServer(dbConn)
	if err := s.Run(*addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
