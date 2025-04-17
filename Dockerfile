# Используй это как финальный образ для разработки (не distroless)
FROM golang:1.24 AS dev

WORKDIR /workspace
COPY . .
RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@latest