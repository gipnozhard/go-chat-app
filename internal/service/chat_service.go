package service

import (
	"errors"
	"strings"

	"go-chat-app/internal/models"
	"go-chat-app/internal/repository"
)

// ChatService содержит бизнес-логику работы с чатами
type ChatService struct {
	chatRepo    *repository.ChatRepository
	messageRepo *repository.MessageRepository
}

// NewChatService создает новый сервис для работы с чатами
func NewChatService(chatRepo *repository.ChatRepository, messageRepo *repository.MessageRepository) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

// CreateChat создает новый чат
func (s *ChatService) CreateChat(title string) (*models.Chat, error) {
	// -------------------------------------------------
	// 1. Триммируем пробелы по краям (как рекомендуется в ТЗ)
	trimmedTitle := strings.TrimSpace(title)
	// -------------------------------------------------

	// 2. Проверяем что title не пустой и длина от 1 до 200
	if len(trimmedTitle) == 0 {
		return nil, errors.New("title не может быть пустым")
	}
	if len(trimmedTitle) > 200 {
		return nil, errors.New("title должен содержать не более 200 символов")
	}

	// 3. Создаем объект чата
	chat := &models.Chat{
		Title: trimmedTitle,
	}

	// 4. Сохраняем в базу
	err := s.chatRepo.Create(chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

// SendMessage отправляет сообщение в чат
func (s *ChatService) SendMessage(chatID uint, text string) (*models.Message, error) {
	// 1. Проверяем что чат существует
	_, err := s.chatRepo.GetByID(chatID)
	if err != nil {
		// Если чат не найден - возвращаем ошибку
		return nil, errors.New("чат не найден")
	}
	// ---------------------------------
	// 2. Триммируем пробелы по краям
	trimmedText := strings.TrimSpace(text)
	// ---------------------------------

	// 3. Проверяем что text не пустой и длина от 1 до 5000
	if len(trimmedText) == 0 {
		return nil, errors.New("текст не может быть пустым")
	}
	if len(trimmedText) > 5000 {
		return nil, errors.New("объем текста должен быть не более 5000 символов")
	}

	// 4. Создаем объект сообщения
	message := &models.Message{
		ChatID: chatID,
		Text:   trimmedText,
	}

	// 5. Сохраняем в базу
	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetChatWithMessages возвращает чат и последние сообщения
func (s *ChatService) GetChatWithMessages(chatID uint, limit int) (*models.Chat, []models.Message, error) {
	// 1. Получаем чат
	chat, err := s.chatRepo.GetByID(chatID)
	if err != nil {
		return nil, nil, errors.New("чат не найден")
	}

	// 2. Ограничиваем limit максимум 100, как в ТЗ
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 20 // значение по умолчанию из ТЗ
	}

	// 3. Получаем последние сообщения
	messages, err := s.messageRepo.GetLastMessagesByChatID(chatID, limit)
	if err != nil {
		return nil, nil, err
	}

	return chat, messages, nil
}

// DeleteChat удаляет чат (сообщения удалятся каскадно через GORM)
func (s *ChatService) DeleteChat(chatID uint) error {
	// 1. Проверяем что чат существует
	_, err := s.chatRepo.GetByID(chatID)
	if err != nil {
		return errors.New("чат не найден")
	}

	// 2. Удаляем чат
	// Сообщения удалятся автоматически благодаря constraint:OnDelete:CASCADE в модели Chat
	return s.chatRepo.Delete(chatID)
}
