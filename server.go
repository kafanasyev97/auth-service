package main

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/kafanasyev97/auth-service/proto/github.com/kafanasyev97/go-microservices-proto/auth"
	// pb "github.com/kafanasyev97/go-microservices/proto/auth"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer

	mu        sync.Mutex
	users     map[string]string // username -> password
	tokens    map[string]string // token -> user_id
	userIDSeq int               // автоинкремент user_id
}

func NewAuthServer() *AuthServer {
	return &AuthServer{
		users:  make(map[string]string),
		tokens: make(map[string]string),
	}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[req.Username]; exists {
		return nil, errors.New("username already exists")
	}

	s.userIDSeq++
	// userId := "user_" + string(rune(s.userIDSeq))
	userId := "user_" + strconv.Itoa(s.userIDSeq)

	s.users[req.Username] = req.Password
	return &auth.RegisterResponse{UserId: userId}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	password, exists := s.users[req.Username]
	if !exists || password != req.Password {
		return nil, errors.New("invalid credentials")
	}

	token := "dummy-token-" + req.Username
	s.tokens[token] = req.Username

	return &auth.LoginResponse{Token: token}, nil

	// простая заглушка
	// return &pb.LoginResponse{
	// 	Token: "dummy-token",
	// }, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.tokens[req.Token]
	if !exists {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	return &auth.ValidateTokenResponse{
		Valid:  true,
		UserId: user,
	}, nil
}
