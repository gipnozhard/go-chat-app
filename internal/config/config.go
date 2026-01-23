package config

import (
	"os"
)

// Config хранит конфигурацию приложения
type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSL  string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		Port:   getEnv("APP_PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "chat_user"),
		DBPass: getEnv("DB_PASSWORD", "chat_password"),
		DBName: getEnv("DB_NAME", "chat_db"),
		DBSSL:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// getEnv получает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
