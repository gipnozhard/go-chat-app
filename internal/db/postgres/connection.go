package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB - глобальная переменная для хранения подключения к базе данных
// Используется во всем приложении для доступа к БД
var DB *gorm.DB

// InitDB инициализирует подключение к базе данных PostgreSQL
// Возвращает *gorm.DB для работы с БД и error в случае неудачи
func InitDB() (*gorm.DB, error) {
	// Загружаем переменные окружения из .env файла
	// Если файла нет - игнорируем ошибку (переменные могут быть уже установлены)
	_ = godotenv.Load()

	// Формируем строку подключения (DSN - Data Source Name) к PostgreSQL
	// Используем переменные окружения для конфигурации:
	// DB_HOST - хост БД (например: localhost или postgres для Docker)
	// DB_USER - имя пользователя БД
	// DB_PASSWORD - пароль пользователя
	// DB_NAME - имя базы данных
	// DB_PORT - порт PostgreSQL (по умолчанию 5432)
	// DB_SSL_MODE - режим SSL (disable для разработки, require для продакшена)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)

	// Подключаемся к базе данных
	var err error
	// Открываем соединение с базой данных используя GORM и PostgreSQL драйвер
	// GORM.Config{} - конфигурация GORM с настройками по умолчанию
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Если подключение не удалось, возвращаем ошибку
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Получаем низкоуровневое соединение *sql.DB из GORM
	// Это нужно для настройки пула соединений
	sqlDB, _ := DB.DB()

	// Настраиваем пул соединений для оптимальной производительности:

	// SetMaxIdleConns - максимальное количество неактивных (idle) соединений в пуле
	// 20 соединений будет поддерживаться в готовности для быстрого использования
	sqlDB.SetMaxIdleConns(20)

	// SetMaxOpenConns - максимальное количество открытых соединений одновременно
	// Ограничивает нагрузку на БД, не позволяет открыть больше 100 соединений
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime - максимальное время жизни соединения
	// Через 1 час соединение будет закрыто и создано новое
	// Это помогает избежать проблем с устаревшими соединениями
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Проверяем соединение с БД
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return DB, nil
}

// GetDB возвращает глобальное подключение к базе данных
// Используется в других частях приложения для получения DB
func GetDB() *gorm.DB {
	return DB
}
