package main

import (
	"fmt"
	"log"
	"net/http"

	"go-chat-app/internal/config"
	"go-chat-app/internal/db/postgres"
	"go-chat-app/internal/db/service"
	"go-chat-app/internal/handler"
	"go-chat-app/internal/repository"
)

func main() {

	// Конфигурация
	cfg := config.LoadConfig()

	// Подключение к БД
	db, err := postgres.InitDB()
	if err != nil {
		log.Fatal("Ошибка БД:", err)
	}

	// ВСЕГДА применяем миграции при запуске
	// goose сам проверяет, какие миграции уже применены
	log.Println("Проверка и применение миграций...")
	if err := postgres.RunMigrations(db, "migrations"); err != nil {
		log.Fatal("Ошибка миграций:", err)
	}
	log.Println("Миграции проверены/применены.")

	// Инициализация зависимостей
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatService := service.NewChatService(chatRepo, messageRepo)
	chatHandler := handler.NewChatHandler(chatService)

	// Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Сервер запущен на http://localhost%s", addr)

	if err := http.ListenAndServe(addr, chatHandler); err != nil {
		log.Fatal("Ошибка сервера:", err)
	}
}
