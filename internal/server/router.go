package server

import (
	"log"
	"net/http"
	"strings"

	"go-chat-app/internal/db/service"
	"go-chat-app/internal/handler"
)

// Router обрабатывает маршрутизацию HTTP запросов
type Router struct {
	chatHandler *handler.ChatHandler
}

// NewRouter создает новый роутер с привязкой хендлеров
func NewRouter(chatService *service.ChatService) *Router {
	return &Router{
		chatHandler: handler.NewChatHandler(chatService),
	}
}

// ServeHTTP обрабатывает все входящие HTTP запросы
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Добавляем middleware
	// 1. Логирование
	// 2. Recovery (обработка паник)
	// 3. Основной обработчик
	handler := r.recoveryMiddleware(r.loggingMiddleware(r.mainHandler))
	handler.ServeHTTP(w, req)
}

// mainHandler определяет какой хендлер вызвать в зависимости от пути
func (r *Router) mainHandler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Простая маршрутизация без сторонних библиотек
	switch {
	// POST /chats/
	case path == "/chats/" && req.Method == http.MethodPost:
		r.chatHandler.CreateChat(w, req)

	// GET /chats/{id}
	case strings.HasPrefix(path, "/chats/") && req.Method == http.MethodGet:
		// Проверяем что после /chats/ идет число (ID чата)
		if isChatIDPath(path) {
			r.chatHandler.GetChat(w, req)
		} else {
			http.NotFound(w, req)
		}

	// DELETE /chats/{id}
	case strings.HasPrefix(path, "/chats/") && req.Method == http.MethodDelete:
		if isChatIDPath(path) {
			r.chatHandler.DeleteChat(w, req)
		} else {
			http.NotFound(w, req)
		}

	// POST /chats/{id}/messages/
	case strings.HasSuffix(path, "/messages/") && req.Method == http.MethodPost:
		if isChatMessagesPath(path) {
			r.chatHandler.SendMessage(w, req)
		} else {
			http.NotFound(w, req)
		}

	// Все остальные пути - 404
	default:
		http.NotFound(w, req)
	}
}

// loggingMiddleware логирует все запросы
func (r *Router) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Простое логирование в консоль
		log.Printf("%s %s from %s", req.Method, req.URL.Path, req.RemoteAddr)
		next(w, req)
	}
}

// recoveryMiddleware ловит паники и возвращает 500 ошибку
func (r *Router) recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Логируем панику
				log.Printf("PANIC: %v", err)
				// Возвращаем 500 ошибку
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, req)
	}
}

// Вспомогательные функции для проверки путей

// isChatIDPath проверяет что путь вида /chats/123 (где 123 - число)
func isChatIDPath(path string) bool {
	// Убираем "/chats/" в начале
	if !strings.HasPrefix(path, "/chats/") {
		return false
	}

	rest := path[len("/chats/"):]
	// Должен остаться только ID (число)
	// Проверяем что в строке только цифры
	for _, char := range rest {
		if char < '0' || char > '9' {
			return false
		}
	}
	return rest != "" // не пустая строка
}

// isChatMessagesPath проверяет что путь вида /chats/123/messages/
func isChatMessagesPath(path string) bool {
	if !strings.HasSuffix(path, "/messages/") {
		return false
	}

	// Убираем "/messages/" в конце
	pathWithoutMessages := path[:len(path)-len("/messages/")]
	// Теперь проверяем что осталось /chats/123
	return isChatIDPath(pathWithoutMessages + "/")
}
