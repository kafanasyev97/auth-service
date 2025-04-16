package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"

	"github.com/kafanasyev97/auth-service/proto/auth"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer

	db        *sql.DB
	mu        sync.Mutex
	users     map[string]string // username -> password
	tokens    map[string]string // token -> user_id
	userIDSeq int               // автоинкремент user_id
}

func NewAuthServer() *AuthServer {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("не удалось подключиться к базе данных: %v", err))
	}

	return &AuthServer{
		db:     db,
		users:  make(map[string]string), // пока оставим для совместимости
		tokens: make(map[string]string),
	}
}

func (s *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	var userID int
	err = s.db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", req.Username, req.Password).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{UserId: fmt.Sprintf("user_%d", userID)}, nil
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
