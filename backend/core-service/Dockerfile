# ---------- builder ----------
FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git

COPY core-service/go.mod core-service/go.sum ./core-service/
COPY core-service ./core-service
# Для клиентских зависимостей (генерённые proto) достаточно самого core-service каталога
RUN cd core-service && go mod download
RUN cd core-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/core ./cmd/main.go

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# положим ДОКЕРНЫЙ конфиг так, как его ищет код (local.yaml)
COPY core-service/config/docker.yaml /app/config/local.yaml
COPY --from=builder /out/core /app/core

EXPOSE 8099
ENTRYPOINT ["/app/core"]