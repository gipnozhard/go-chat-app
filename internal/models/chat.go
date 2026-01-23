package models

import (
	"time"

	"gorm.io/gorm"
)

// Chat представляет собой модель чата в базе данных
// GORM автоматически создаст таблицу "chats" на основе этой структуры
type Chat struct {
	// ID - первичный ключ, автоинкремент
	// gorm:"primaryKey" - указывает GORM что это первичный ключ
	// json:"id" - при сериализации в JSON поле будет называться "id"
	ID uint `gorm:"primaryKey" json:"id"`

	// Title - заголовок чата, обязательное поле
	// gorm:"size:200;not null" - ограничения в БД:
	//   size:200 - максимальная длина 200 символов (VARCHAR(200))
	//   not null - поле не может быть NULL
	// json:"title" - в JSON будет как "title"
	Title string `gorm:"size:200;not null" json:"title"`

	// Временные метки

	// CreatedAt - время создания записи
	// GORM автоматически заполняет это поле при создании
	// json:"created_at" - в JSON будет в формате ISO 8601
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt - время последнего обновления
	// GORM автоматически обновляет это поле при изменении
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt - время "мягкого" удаления (soft delete)
	// gorm:"index" - создает индекс для ускорения поиска удаленных записей
	// json:"-" - НЕ включать это поле в JSON ответы
	// gorm.DeletedAt - специальный тип GORM для soft delete
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Messages - связь "один ко многим" с сообщениями
	// []Message - слайс (массив) сообщений
	// gorm:"foreignKey:ChatID" - указывает какое поле в Message ссылается на Chat
	// constraint:OnDelete:CASCADE - при удалении чата ВСЕ его сообщения удалятся автоматически
	// json:"messages,omitempty" - в JSON будет как "messages", если не пустой
	Messages []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE;" json:"messages,omitempty"`
}
