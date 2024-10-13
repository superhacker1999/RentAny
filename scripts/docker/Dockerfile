# Используем официальный образ Golang
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем все файлы проекта
COPY cmd ./cmd

COPY internal ./internal

COPY .env .

# Сборка бинарника Go
RUN go build -o server ./cmd/main.go

# Используем легковесный образ для запуска приложения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарник из предыдущего этапа
COPY --from=builder /app/server .

COPY --from=builder /app/.env .

# Открываем порт 8080 (или тот, который у вас используется)
EXPOSE 8080

# Команда для запуска сервера
CMD ["./server"]