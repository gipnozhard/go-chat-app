package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheck проверяет, что сервер корректно обрабатывает health check запрос
// Health check - это endpoint для проверки работоспособности приложения
func TestHealthCheck(t *testing.T) {
	// Создаем экземпляр ChatHandler для тестирования
	// nil передается как сервис, потому что для /health endpoint
	// не требуется бизнес-логика (он не зависит от сервиса)
	handler := NewChatHandler(nil)

	// httptest.NewRequest создает фиктивный HTTP запрос
	// Параметры:
	//   "GET" - HTTP метод
	//   "/health" - URL путь (health check endpoint)
	//   nil - тело запроса (пустое, так как GET запрос)
	req := httptest.NewRequest("GET", "/health", nil)

	// httptest.NewRecorder имитирует ResponseWriter
	// Он записывает все что отправляет handler: статус код, заголовки, тело
	// Позволяет проверить ответ без запуска реального сервера
	rr := httptest.NewRecorder()

	// Вызываем метод ServeHTTP, который обрабатывает ВСЕ HTTP запросы
	// Внутри него есть логика маршрутизации, которая определит что
	// путь "/health" должен обработаться health check
	handler.ServeHTTP(rr, req)

	// Health check должен всегда возвращать 200 OK если сервер работает
	// Это стандарт для health check endpoints
	if rr.Code != http.StatusOK {
		// t.Errorf выводит ошибку если тест не прошел
		// Сообщение поможет понять что пошло не так
		t.Errorf("Ожидался статус 200 (OK), получен %d", rr.Code)
	}

	// Ожидаемое тело ответа - JSON с ключом "status" и значением "ok"
	// Это стандартный формат для health check
	expected := `{"status":"ok"}`

	// Сравниваем фактическое тело ответа с ожидаемым
	if rr.Body.String() != expected {
		t.Errorf(
			"Ожидалось тело ответа: %s\nПолучено: %s",
			expected,
			rr.Body.String(),
		)
	}
}

//   go test ./internal/handler -v -run TestHealthCheck
//
// Результат при успехе:
//   === RUN   TestHealthCheck
//   --- PASS: TestHealthCheck (0.00s)
//   PASS
//
// Результат при ошибке (пример):
//   === RUN   TestHealthCheck
//   --- FAIL: TestHealthCheck (0.00s)
//   handler_test.go:XX: Ожидался статус 200, получен 404
