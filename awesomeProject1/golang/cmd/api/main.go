package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang/internal/handlers"
	"golang/internal/middleware"
	"golang/internal/repository"
	"golang/internal/repository/_postgres"
	"golang/internal/usecase"
	"golang/pkg/modules"
)

func main() {
	cfg := &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "Birzhan1234;", // <-- твой рабочий пароль
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}

	ctx := context.Background()
	pg := _postgres.NewPGXDialect(ctx, cfg)

	repos := repository.NewRepositories(pg)
	userUC := usecase.NewUserUsecase(repos.UserRepository)
	h := handlers.NewHandler(userUC)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/users", h.Users)
	mux.HandleFunc("/users/", h.UserByID)

	var handler http.Handler = mux
	handler = middleware.APIKey("secret12345")(handler)
	handler = middleware.Logging(handler)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
