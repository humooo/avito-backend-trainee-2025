# Сборка (используем 1.23, чтобы совпадало с go.mod)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Качаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код и миграции
COPY . .

# Собираем бинарник
RUN go build -o main ./cmd/app/main.go

# Запуск (минимальный образ)
FROM alpine:latest
WORKDIR /root/

# Копируем бинарник из прошлого этапа
COPY --from=builder /app/main .
# Копируем миграции, чтобы main.go мог их прочитать
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
