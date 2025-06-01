# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/app .

# Копируем шаблоны
COPY --from=builder /app/templates ./templates

# Запускаем приложение
CMD ["./app"] 