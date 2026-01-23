package models

import (
	"time"
)

// Message представляет собой модель сообщения в чате
type Message struct {
	// ID - уникальный идентификатор сообщения
	ID uint `gorm:"primaryKey" json:"id"`

	// ChatID - внешний ключ на таблицу чатов
	// not null - сообщение всегда должно быть привязано к чату
	// index - индекс для быстрого поиска сообщений по chat_id
	ChatID uint `gorm:"not null;index" json:"chat_id"`

	// Text - текст сообщения
	// type:text - поле TEXT в БД (поддерживает длинные сообщения до 5000 символов)
	// not null - сообщение не может быть пустым
	Text string `gorm:"type:text;not null" json:"text"`

	// Временные метки, ОПИСАННИЕ МОЖНО ПОСМОТРЕТЬ models/chat.go
	CreatedAt time.Time `json:"created_at"`
}
