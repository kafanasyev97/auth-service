# Используй это как финальный образ для разработки (не distroless)
FROM golang:1.24 AS dev

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o auth-service .
