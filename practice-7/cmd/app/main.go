package main

import (
	"log"
	"practice-7/internal/controller/http/v1"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/postgres"
	"practice-7/utils" 

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Предупреждение: .env файл не найден, используются системные переменные")
	}

	pg, err := postgres.New()
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}

	log.Println("Запуск миграции базы данных...")
	err = pg.Conn.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}
	log.Println("Миграция завершена успешно!")

	userRepo := repo.NewUserRepo(pg.Conn)
	userUseCase := usecase.NewUserUseCase(userRepo)

	handler := gin.Default()

	handler.Use(utils.RateLimiterMiddleware())

	v1Group := handler.Group("/v1")
	v1.NewUserRoutes(v1Group, userUseCase)

	log.Println("Сервер успешно запущен на порту :8090")
	if err := handler.Run(":8090"); err != nil {
		log.Fatal("Не удалось запустить сервер: ", err)
	}
}