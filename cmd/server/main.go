package main

import (
	"log"
	"net"
	"os"

	"github.com/kafanasyev97/auth-service/internal/handler"
	"github.com/kafanasyev97/auth-service/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("не удалось слушать порт: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	redisHost := os.Getenv("REDIS_HOST")
	log.Printf("Подключаемся к БД на %s и Redis на %s\n", dbHost, redisHost)

	grpcServer := grpc.NewServer()
	authHandler := handler.NewAuthHandler()
	auth.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println("Auth Service запущен на порту 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}
