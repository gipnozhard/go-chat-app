# Сборка
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/

# Заключительный этап
FROM alpine:latest
WORKDIR /app

# Копируем из билдера
COPY --from=builder /app/main .
# Копируем папку миграций
COPY --from=builder /app/migrations ./migrations

# Указание порта
EXPOSE 8080

# Команда запуска
CMD ["./main"]