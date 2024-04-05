package main

import (
	"log"
	"os"

	"github.com/Cr4z1k/Avito-test-task/internal/db"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
	"github.com/Cr4z1k/Avito-test-task/internal/service"
	"github.com/Cr4z1k/Avito-test-task/internal/transport/rest"
	"github.com/Cr4z1k/Avito-test-task/internal/transport/rest/handlers"
	"github.com/Cr4z1k/Avito-test-task/pkg/auth"
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

	tokenManager, err := auth.NewTokenManager(os.Getenv("JWT_SALT"))
	if err != nil {
		log.Fatal("Fatal creating tokenManager: ", err.Error())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo, tokenManager)
	handler := handlers.NewHandler(service)

	if err := s.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatal("Fatal start server: ", err.Error())
	}
}
