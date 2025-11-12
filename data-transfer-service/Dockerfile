# ---------- builder ----------
FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git

# моды (build context — корень)
COPY data-transfer-service/go.mod data-transfer-service/go.sum ./data-transfer-service/
COPY core-service/go.mod core-service/go.sum ./core-service/

RUN cd data-transfer-service && go mod download

# код
COPY data-transfer-service ./data-transfer-service
COPY core-service ./core-service

RUN cd data-transfer-service && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/data-transfer ./cmd/main.go

# ---------- runtime ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/data-transfer /app/data-transfer

EXPOSE 50052
ENTRYPOINT ["/app/data-transfer"]