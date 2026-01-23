package postgres

import (
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// RunMigrations применяет миграции базы данных
// migrationsDir - путь к папке с файлами миграций
func RunMigrations(db *gorm.DB, migrationsDir string) error {
	// Получаем низкоуровневое соединение *sql.DB из GORM
	// goose работает со стандартным sql.DB, а не с GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("не удалось получить sql.DB: %w", err)
	}

	// Устанавливаем диалект PostgreSQL
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("не удалось установить диалект БД: %w", err)
	}

	// Применяем миграции
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("не удалось применить миграции: %w", err)
	}

	log.Println("Миграции базы данных успешно применены!")
	return nil
}
