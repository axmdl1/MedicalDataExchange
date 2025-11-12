# ---------- builder ----------
FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git

# openssl не нужен, но git обязателен для private deps
RUN apk add --no-cache git

# скопируем только модули, чтобы кешировать go mod download
COPY user-service/go.mod user-service/go.sum ./user-service/
# для replace => ../core-service / ../data-transfer-service нужен весь контекст,
# но build context у нас корень, так что всё доступно
COPY core-service/go.mod core-service/go.sum ./core-service/
COPY data-transfer-service/go.mod data-transfer-service/go.sum ./data-transfer-service/

RUN cd user-service && go mod download

# теперь код
COPY user-service ./user-service
COPY core-service ./core-service
COPY data-transfer-service ./data-transfer-service

# соберём бинарь
RUN cd user-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/user-service ./cmd/main.go

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# конфиг для докера
COPY user-service/config/docker.yaml /app/config/docker.yaml
COPY --from=builder /out/user-service /app/user-service

EXPOSE 50053
ENTRYPOINT ["/app/user-service", "--config=/app/config/docker.yaml"]