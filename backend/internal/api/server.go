package api

import (
	"context"
	"errors"
	"example/internal/services/mail"
	"net/http"

	"example/internal/database"

	"github.com/minio/minio-go/v7"
)

type Message string

type Response struct {
	Message     Message
	RedirectURL string
}

var (
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
)

type Config struct {
	DB     *database.DB
	MinIO  *minio.Client
	Mailer *mail.Mailer
}

func NewServer(cfg Config) *Config {
	return &Config{DB: cfg.DB, MinIO: cfg.MinIO, Mailer: cfg.Mailer}
}

func (s *Config) RootRoute(ctx context.Context, r *http.Request) ([]byte, error) {
	return []byte("hey, what's that over there?!"), nil
}

func (s *Config) Healthz(ctx context.Context, r *http.Request) ([]byte, error) {
	return []byte("OK"), nil
}
