package main

import (
	"log"
	"net"

	// pb "github.com/kafanasyev97/go-microservices/proto/auth"
	// pb "github.com/kafanasyev97/go-microservices-proto/auth"
	"github.com/kafanasyev97/auth-service/proto/github.com/kafanasyev97/go-microservices-proto/auth"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("не удалось слушать порт: %v", err)
	}

	grpcServer := grpc.NewServer()
	authServer := NewAuthServer()
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	log.Println("Auth Service запущен на порту 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("не удалось запустить сервер: %v", err)
	}
}
