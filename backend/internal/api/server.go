package api

import (
	"errors"
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
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")
)

type Config struct {
	DB    *database.DB
	MinIO *minio.Client
}

func NewServer(cfg Config) *Config {
	return &Config{DB: cfg.DB, MinIO: cfg.MinIO}
}

func (s *Config) HandleRootRoute(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hi :)"))
}
