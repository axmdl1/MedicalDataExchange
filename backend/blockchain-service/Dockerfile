FROM golang:1.23-alpine AS builder

WORKDIR /build

# Копируем go.mod и go.sum для всех зависимых сервисов
COPY blockchain-service/go.mod blockchain-service/go.sum ./blockchain-service/
COPY core-service/go.mod core-service/go.sum ./core-service/

# Копируем сгенерированные proto файлы из core-service
COPY core-service/pkg ./core-service/pkg

# Скачиваем зависимости
RUN cd blockchain-service && GOTOOLCHAIN=auto go mod download

# Копируем исходники blockchain-service
COPY blockchain-service/ ./blockchain-service/

# Собираем приложение
RUN cd blockchain-service && GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux go build -o blockchain-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем бинарник
COPY --from=builder /build/blockchain-service/blockchain-service .

# Создаем директории для конфигурации
RUN mkdir -p /app/config /app/wallet /app/crypto

EXPOSE 50054

CMD ["./blockchain-service"]
