package main

import (
	"context"
	"errors"
	"example/internal/api"
	"example/internal/database"
	"example/internal/middleware"
	"example/internal/services/mail"
	"example/internal/services/minio"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var port = 1323

func main() {
	// Init database
	conn := database.NewClient()
	defer conn.DB.Close()

	// Init minio
	minioClient, err := minio.New(minio.Config{
		Host:      os.Getenv("MINIO_HOST"),
		Port:      os.Getenv("MINIO_PORT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("MINIO_SECRET_ACCESS_KEY"),
		SSL:       os.Getenv("MINIO_SSL") == "true",
	})

	// Mailer
	mailer := mail.NewClient()

	// Init router
	router := http.NewServeMux()
	handler := api.NewServer(api.Config{
		DB:     conn,
		MinIO:  minioClient,
		Mailer: &mailer,
	})

	// Middlewares
	stack := middleware.CreateStack(
		middleware.CORS,
		middleware.Authentication,
	)

	// Public routes
	router.HandleFunc("GET /", wrap(handler.RootRoute))
	router.HandleFunc("GET /healthz", wrap(handler.Healthz))

	// Auth routes
	router.HandleFunc("POST /one_time_login", wrap(handler.OneTimeLogin))
	router.HandleFunc("POST /sign_in", wrap(handler.SignIn))
	router.HandleFunc("POST /sign_up", wrap(handler.SignUp))
	router.HandleFunc("POST /logout", wrap(handler.LogOut))

	// File routes
	router.HandleFunc("GET /files", wrap(handler.Files))
	router.HandleFunc("POST /files", wrap(handler.FileUpload))

	router.HandleFunc("PATCH /files/{id}", wrap(handler.FilePatch))
	router.HandleFunc("DELETE /files/{id}", wrap(handler.FileDelete))

	router.HandleFunc("GET /files/{id}/preview", wrap(handler.FilePreview))
	router.HandleFunc("GET /files/{id}/download", handler.FileDownload)

	// Folder routes
	router.HandleFunc("GET /folders/{id}", wrap(handler.Folders))

	// Shared drive routes
	router.HandleFunc("GET /shared_drives", wrap(handler.SharedDrives))

	// Server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: stack(router),
	}

	slog.Info(fmt.Sprintf("starting server on http://localhost:%d", port))

	// Start server
	err = server.ListenAndServe()
	if err != nil {
		slog.Error("error starting server", "err", err)
	}
}

// wrap is a helper function to wrap the http handler functions with error handling
func wrap(handler func(ctx context.Context, r *http.Request) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := handler(r.Context(), r)

		switch {
		case errors.Is(err, api.ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
			return
		case errors.Is(err, api.ErrBadRequest):
			w.WriteHeader(http.StatusBadRequest)
			return
		case errors.Is(err, api.ErrInternal):
			w.WriteHeader(http.StatusInternalServerError)
			return
		case err != nil:
			slog.Error("error handling request", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(res)
			if err != nil {
				return
			}
		}
	}
}
