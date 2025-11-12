# ---------- builder ----------
FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git

COPY user-service/go.mod user-service/go.sum ./user-service/
COPY core-service/go.mod core-service/go.sum ./core-service/
COPY data-transfer-service/go.mod data-transfer-service/go.sum ./data-transfer-service/
RUN cd user-service && go mod download

COPY user-service ./user-service
COPY core-service ./core-service
COPY data-transfer-service ./data-transfer-service
RUN cd user-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/user-service ./cmd/main.go

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# ключевая строка: кладём конфиг туда, где его ищет код (config/local.yaml)
COPY user-service/config/docker.yaml /app/config/local.yaml
COPY --from=builder /out/user-service /app/user-service

EXPOSE 50053
ENTRYPOINT ["/app/user-service"]