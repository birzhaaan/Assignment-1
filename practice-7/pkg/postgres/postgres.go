package postgres

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	Conn *gorm.DB
}

func New() (*Postgres, error) {
	dbPass := os.Getenv("DB_PASSWORD") 
	
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=practice_db port=5432 sslmode=disable", dbPass)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Postgres{Conn: db}, nil
}