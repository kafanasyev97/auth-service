package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kafanasyev97/auth-service/proto/auth"
	_ "github.com/lib/pq"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer

	db     *sql.DB
	mu     sync.Mutex
	tokens map[string]string
}

func NewAuthHandler() *AuthHandler {
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

	return &AuthHandler{
		db:     db,
		tokens: make(map[string]string),
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	var userID int
	err = h.db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", req.Username, req.Password).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{UserId: fmt.Sprintf("user_%d", userID)}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var userID int
	var password string
	err := h.db.QueryRow("SELECT id, password FROM users WHERE username = $1", req.Username).Scan(&userID, &password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if password != req.Password {
		return nil, errors.New("invalid password")
	}

	token := fmt.Sprintf("token_%d", time.Now().UnixNano())
	h.tokens[token] = fmt.Sprintf("user_%d", userID)

	return &auth.LoginResponse{Token: token}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userID, exists := h.tokens[req.Token]
	if !exists {
		return &auth.ValidateTokenResponse{Valid: false}, nil
	}

	return &auth.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
