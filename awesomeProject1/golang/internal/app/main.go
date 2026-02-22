package app

import (
	"context"
	"fmt"
	"time"

	"golang/internal/repository"
	"golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := initPostgreConfig()
	pg := _postgres.NewPGXDialect(ctx, dbConfig)

	repos := repository.NewRepositories(pg)

	users, err := repos.GetUsers()
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}

	fmt.Printf("Users: %+v\n", users)
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "Birzhan1234;", // оставь как у тебя рабочий
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
