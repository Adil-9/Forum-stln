package main

import (
	"flag"
	"fmt"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	handler, server, err := manage()
	if err != nil {
		return
	}

	fmt.Printf("Starting server at addr %s\nhttp://localhost%s/\n", *addr, *addr)

	if err := server.Run(*addr, handler.InitRoutes()); err != nil {
		log.Fatalf("error while running the server: %s", err.Error())
	}
}

func manage() (*delivery.Handler, *server.Server, error) {
	db, err := repository.OpenSqliteDB("store.db")
	if err != nil {
		log.Fatalf("error while opening db: %s", err)
		return nil, nil, err
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := delivery.NewHandler(service)
	server := new(server.Server)

	return handler, server, nil
}
