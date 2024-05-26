package main

import (
	"example/internal/api"
	"example/internal/database"
	"example/internal/middleware"
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

	// Init router
	router := http.NewServeMux()
	handler := api.NewServer(api.Config{
		DB:    conn,
		MinIO: minioClient,
	})

	router.HandleFunc("/", handler.HandleRootRoute)

	//Authentication
	router.HandleFunc("POST /sign_in", handler.HandleSignIn)

	// User routes
	//router.HandleFunc("GET /me", handler.HandleMe)

	stack := middleware.CreateStack(
		middleware.CORS,
	)

	// File routes
	router.HandleFunc("DELETE /files/{id}", handler.HandleFileDelete)
	router.HandleFunc("PATCH /files/{id}", handler.HandleFilePatch)
	router.HandleFunc("GET /files/{id}/download", handler.HandleFileDownload)
	router.HandleFunc("GET /files", handler.HandleFiles)
	router.HandleFunc("POST /files", handler.HandleFileUpload)
	router.HandleFunc("GET /folders/{id}", handler.HandleFolders)
	router.HandleFunc("GET /shared-drives", handler.HandleSharedDrives)

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
