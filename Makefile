.PHONY: help build run clean

help:
	@echo "Доступные команды:"
	@echo "  make build    - Создание образов Docker"
	@echo "  make run      - Запуск приложения"
	@echo "  make clean    - Остановка и удаление контейнеры"
	@echo "  make test     - Запуск тестов"
	@echo "  make          - Показать эту справку"

build:
	@echo "Создание образов Docker..."
	docker-compose build

run:
	@echo "Запуск приложения..."
	docker-compose up

clean:
	@echo "Отчистка..."
	docker-compose down -v

test:
	@echo "Запуск тестов..."
	go test ./internal/handler -v -run TestHealthCheck