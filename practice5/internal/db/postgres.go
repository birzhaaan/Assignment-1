package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=Birzhan1234; dbname=practice5 sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return database
}
