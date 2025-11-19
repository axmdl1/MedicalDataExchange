FROM golang:1.25

WORKDIR /app

# git нужен для go mod download
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o deploy ./cmd/demo/main.go

CMD ["/app/deploy"]