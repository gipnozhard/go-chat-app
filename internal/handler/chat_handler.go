package handler

import (
	"encoding/json"
	"go-chat-app/internal/db/service"
	"net/http"
	"strconv"
	"strings"

	"go-chat-app/internal/models"
)

// ChatHandler обрабатывает HTTP запросы для работы с чатами и сообщениями
type ChatHandler struct {
	service *service.ChatService // Сервис с бизнес-логикой
}

// ServeHTTP обрабатывает все входящие HTTP запросы и перенаправляет их на соответствующие методы
// Этот метод реализует интерфейс http.Handler
func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Простая маршрутизация на основе пути и метода запроса

	// Проверяем все возможные комбинации пути и метода:
	switch {
	// СЛУЧАЙ 1: Создание нового чата
	// Путь: POST /chats
	// Пример: POST http://localhost:8080/chats
	case r.URL.Path == "/chats" && r.Method == "POST":
		h.CreateChat(w, r)

	// СЛУЧАЙ 2: Отправка сообщения в чат
	// Путь: POST /chats/{id}/messages
	// Пример: POST http://localhost:8080/chats/123/messages
	case strings.HasPrefix(r.URL.Path, "/chats/") && strings.HasSuffix(r.URL.Path, "/messages") && r.Method == "POST":
		h.SendMessage(w, r)

	// СЛУЧАЙ 3: Получение информации о чате с сообщениями
	// Путь: GET /chats/{id}
	// Пример: GET http://localhost:8080/chats/123?limit=20

	case strings.HasPrefix(r.URL.Path, "/chats/") && r.Method == "GET":
		h.GetChat(w, r)

	// СЛУЧАЙ 4: Удаление чата
	// Путь: DELETE /chats/{id}
	// Пример: DELETE http://localhost:8080/chats/123
	case strings.HasPrefix(r.URL.Path, "/chats/") && r.Method == "DELETE":
		h.DeleteChat(w, r)

	// ВАРИАНТ 5: HEALTH CHECK (для мониторинга)
	// Условие: путь "/health" И метод GET
	// Используется Docker, Kubernetes и т.д. для проверки что сервер жив
	case r.URL.Path == "/health" && r.Method == "GET":
		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// Пишем простой JSON ответ
		w.Write([]byte(`{"status":"ok"}`))
	default:

		// ВАРИАНТ 6: НЕИЗВЕСТНЫЙ ПУТЬ
		// Если ни одно из условий выше не выполнилось - путь не существует
		// Возвращаем стандартную 404 ошибку "Not Found"
		http.NotFound(w, r)
	}
}

// NewChatHandler создает новый обработчик чатов
func NewChatHandler(service *service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

// 1. POST /chats/ - создать новый чат
// Тело запроса: {"title": "Название чата"}
// Ответ: созданный чат в формате JSON
func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что используется правильный HTTP метод
	if r.Method != "POST" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed) // 405
		return
	}

	// Структура для парсинга JSON тела запроса
	var data struct {
		Title string `json:"title"` // Название чата
	}

	// Декодируем JSON тело запроса
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest) // 400
		return
	}

	// Вызываем сервис для создания чата
	chat, err := h.service.CreateChat(data.Title)
	if err != nil {
		// Обрабатываем ошибки валидации (400) и остальные (500)
		if strings.Contains(err.Error(), "не может быть пустым") ||
			strings.Contains(err.Error(), "не более") {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
		} else {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError) // 500
		}
		return
	}

	// Успешный ответ: возвращаем созданный чат
	// 1. Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	// Теперь клиент знает, что получит JSON

	// 2. Устанавливаем статус код 201 Created
	w.WriteHeader(http.StatusCreated) // 201 Created
	// Это стандартный статус для успешно созданного ресурса

	// 3. Кодируем объект chat в JSON и отправляем
	json.NewEncoder(w).Encode(chat)
	// GORM автоматически заполнил chat.ID, chat.CreatedAt и т.д.
}

// 2. POST /chats/{id}/messages/ - отправить сообщение в чат
// Тело запроса: {"text": "Текст сообщения"}
// Ответ: созданное сообщение в формате JSON
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Проверяем HTTP метод
	if r.Method != "POST" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed) // 405
		return
	}

	// Разбираем URL путь для получения ID чата
	// Пример: /chats/123/messages/ → parts = ["chats", "123", "messages"]
	path := strings.Trim(r.URL.Path, "/") // Убираем слэши в начале и конце
	parts := strings.Split(path, "/")     // Разбиваем по слэшам

	// Проверяем структуру пути: должно быть 3 части
	if len(parts) != 3 || parts[0] != "chats" || parts[2] != "messages" {
		http.Error(w, "Неверный URL", http.StatusBadRequest) // 400
		return
	}

	// Преобразуем ID чата из строки в число
	chatID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Неверный ID чата", http.StatusBadRequest) // 400
		return
	}

	// Структура для парсинга JSON тела запроса
	var data struct {
		Text string `json:"text"` // Текст сообщения
	}

	// Декодируем JSON тело запроса
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest) // 400
		return
	}

	// Вызываем сервис для отправки сообщения
	message, err := h.service.SendMessage(uint(chatID), data.Text)
	if err != nil {
		// Разные типы ошибок = разные HTTP статусы
		if strings.Contains(err.Error(), "не найден") {
			http.Error(w, "Чат не найден", http.StatusNotFound) // 404
		} else if strings.Contains(err.Error(), "не может быть пустым") ||
			strings.Contains(err.Error(), "не более") {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
		} else {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError) // 500
		}
		return
	}

	// Успешный ответ: возвращаем созданное сообщение
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(message)
}

// 3. GET /chats/{id} - получить информацию о чате и его сообщениях
// Query параметр: limit (по умолчанию 20, максимум 100)
// Ответ: {"chat": {...}, "messages": [...]}
func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	// Проверяем HTTP метод
	if r.Method != "GET" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed) // 405
		return
	}

	// Разбираем URL путь для получения ID чата
	// Пример: /chats/123 → parts = ["chats", "123"]
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// Проверяем структуру пути: должно быть 2 части
	if len(parts) != 2 || parts[0] != "chats" {
		http.Error(w, "Неверный URL", http.StatusBadRequest) // 400
		return
	}

	// Преобразуем ID чата из строки в число
	chatID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Неверный ID чата", http.StatusBadRequest) // 400
		return
	}

	// Получаем параметр limit из query строки
	// Пример: /chats/123?limit=50
	limit := 20 // Значение по умолчанию из ТЗ
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		// Пытаемся преобразовать строку в число
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			// Ограничиваем максимум 100 сообщений (требование ТЗ)
			if limit > 100 {
				limit = 100
			}
		}
	}

	// Вызываем сервис для получения чата и сообщений
	chat, messages, err := h.service.GetChatWithMessages(uint(chatID), limit)
	if err != nil {
		// Обрабатываем ошибки
		if strings.Contains(err.Error(), "не найден") {
			http.Error(w, "Чат не найден", http.StatusNotFound) // 404
		} else {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError) // 500
		}
		return
	}

	// Формируем и возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	// Анонимная структура для ответа
	json.NewEncoder(w).Encode(struct {
		Chat     models.Chat      `json:"chat"`     // Информация о чате
		Messages []models.Message `json:"messages"` // Список сообщений
	}{
		Chat:     *chat,
		Messages: messages,
	})
}

// 4. DELETE /chats/{id} - удалить чат и все его сообщения
// Ответ: 204 No Content
func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	// Проверяем HTTP метод
	if r.Method != "DELETE" {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed) // 405
		return
	}

	// Разбираем URL путь для получения ID чата
	// Пример: /chats/123 → parts = ["chats", "123"]
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	// Проверяем структуру пути: должно быть 2 части
	if len(parts) != 2 || parts[0] != "chats" {
		http.Error(w, "Неверный URL", http.StatusBadRequest) // 400
		return
	}

	// Преобразуем ID чата из строки в число
	chatID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Неверный ID чата", http.StatusBadRequest) // 400
		return
	}

	// Вызываем сервис для удаления чата
	err = h.service.DeleteChat(uint(chatID))
	if err != nil {
		// Обрабатываем ошибки
		if strings.Contains(err.Error(), "не найден") {
			http.Error(w, "Чат не найден", http.StatusNotFound) // 404
		} else {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError) // 500
		}
		return
	}

	// Успешный ответ: 204 No Content (как указано в ТЗ)
	w.WriteHeader(http.StatusNoContent) // 204
}
