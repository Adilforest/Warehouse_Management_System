# Этап сборки
FROM golang:1.23.4 AS builder

WORKDIR /app

# Копируем исходный код и файл .env
COPY . .

# Скачиваем зависимости и собираем приложение
RUN go mod download
RUN go build -o main .

# Этап финального образа
FROM gcr.io/distroless/base

WORKDIR /app

# Копируем собранное приложение и файл .env
COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /app/.env  

# Запускаем приложение
CMD ["/app/main"]