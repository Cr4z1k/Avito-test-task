package main

import (
	"log"

	"github.com/Cr4z1k/Avito-test-task/internal/db"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
	"github.com/Cr4z1k/Avito-test-task/internal/service"
	"github.com/Cr4z1k/Avito-test-task/internal/transport/rest"
	"github.com/Cr4z1k/Avito-test-task/internal/transport/rest/handlers"
	"github.com/joho/godotenv"
)

func main() {
	s := new(rest.Server)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Fatal loading .env: ", err.Error())
	}

	db, err := db.GetConnection()
	if err != nil {
		log.Fatal("Fatal connect to DB: ", err.Error())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handlers.NewHandler(service)

	if err := s.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatal("Fatal start server: ", err.Error())
	}
}
