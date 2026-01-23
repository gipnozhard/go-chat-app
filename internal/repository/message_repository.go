package repository

import (
	"go-chat-app/internal/models"

	"gorm.io/gorm"
)

// MessageRepository отвечает за работу с сообщениями в базе данных
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository создает новый репозиторий для сообщений
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create сохраняет новое сообщение в базу данных
func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// GetLastMessagesByChatID возвращает последние сообщения чата
// limit - сколько сообщений вернуть, отсортированные по created_at (новые первые)
func (r *MessageRepository) GetLastMessagesByChatID(chatID uint, limit int) ([]models.Message, error) {
	var messages []models.Message

	// Where - фильтр по chat_id
	// Limit - ограничение количества
	err := r.db.Where("chat_id = ?", chatID).
		Limit(limit).
		Find(&messages).Error

	return messages, err
}
