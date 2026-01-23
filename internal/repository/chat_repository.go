package repository

import (
	"go-chat-app/internal/models"

	"gorm.io/gorm"
)

// ChatRepository отвечает за работу с чатами в базе данных
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository создает новый репозиторий для чатов
func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// Create сохраняет новый чат в базу данных
func (r *ChatRepository) Create(chat *models.Chat) error {
	return r.db.Create(chat).Error
}

// GetByID находит чат по ID
func (r *ChatRepository) GetByID(id uint) (*models.Chat, error) {
	var chat models.Chat
	// First ищет первую запись по условию
	err := r.db.First(&chat, id).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// Delete удаляет чат по ID
func (r *ChatRepository) Delete(id uint) error {
	// Delete удаляет запись по ID
	return r.db.Delete(&models.Chat{}, id).Error
}
