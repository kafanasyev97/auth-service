FROM golang:1.24 as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o auth-service

FROM gcr.io/distroless/base-debian12

WORKDIR /
COPY --from=builder /app/auth-service /auth-service

ENTRYPOINT ["/auth-service"]
