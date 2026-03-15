package main

import (
	"github.com/gin-gonic/gin"
	"practice5/internal/db"
	"practice5/internal/handler"
	"practice5/internal/repository"
)

func main() {
	database := db.NewPostgresDB()
	defer database.Close()

	repo := repository.NewRepository(database)
	h := handler.NewHandler(repo)

	r := gin.Default()

	r.GET("/users", h.GetUsers)
	r.GET("/common-friends", h.GetCommonFriends)

	r.Run(":8080")
}
