# 1
# Используем официальный образ Golang для сборки приложения
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем файл go.mod
COPY go.mod ./
COPY go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем приложение. Предполагается, что точка входа находится в cmd/server/main.go
RUN CGO_ENABLED=0 go build -o app ./cmd/main.go

# 2
# минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем собранное приложение из первого этапа
COPY --from=builder /app/app .

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["/app/app"]
