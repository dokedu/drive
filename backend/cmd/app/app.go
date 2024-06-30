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

	stack := middleware.CreateStack(
		middleware.CORS,
		middleware.Authentication,
	)

	// Public routes
	router.HandleFunc("GET /", handler.HandleRootRoute)
	router.HandleFunc("POST /one_time_login", wrapper(handler.HandleOneTimeLogin))
	router.HandleFunc("POST /login", wrapper(handler.HandleLogin))
	router.HandleFunc("POST /sign_up", wrapper(handler.HandleSignUp))

	router.HandleFunc("POST /logout", wrapper(handler.HandleLogOut))

	// File routes
	router.HandleFunc("DELETE /files/{id}", wrapper(handler.HandleFileDelete))
	router.HandleFunc("PATCH /files/{id}", wrapper(handler.HandleFilePatch))
	router.HandleFunc("GET /files/{id}/download", handler.HandleFileDownload)
	router.HandleFunc("GET /files/{id}/preview", wrapper(handler.HandleFilePreview))
	router.HandleFunc("GET /files", wrapper(handler.HandleFiles))
	router.HandleFunc("POST /files", wrapper(handler.HandleFileUpload))
	router.HandleFunc("GET /folders/{id}", wrapper(handler.HandleFolders))
	router.HandleFunc("GET /shared_drives", wrapper(handler.HandleSharedDrives))

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

func wrapper(handler func(ctx context.Context, r *http.Request) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := handler(ctx, r)
		if err != nil {
			if errors.Is(err, api.ErrUnauthorized) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if errors.Is(err, api.ErrBadRequest) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			slog.Error("error handling request", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(res)
		if err != nil {
			slog.Error("error writing response", "err", err)
			return
		}
		return
	}
}
