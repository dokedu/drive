package main

import (
	"example/internal/api"
	"example/internal/database"
	"example/internal/middleware"
	"example/internal/services/mail"
	"example/internal/services/minio"
	"fmt"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	handler := api.NewServer(api.Config{
		DB:     conn,
		MinIO:  minioClient,
		Mailer: &mailer,
	})

	// Public routes
	router.HandleFunc("/", handler.HandleRootRoute).Methods("GET")
	router.HandleFunc("/one-time-login", handler.HandleOneTimeLogin).Methods("POST")
	router.HandleFunc("/login", handler.HandleLogin).Methods("POST")

	// File routes
	fileRouter := router.PathPrefix("/files/").Subrouter()
	fileRouter.HandleFunc("/{id}", handler.HandleFileDelete).Methods("DELETE")
	fileRouter.HandleFunc("/{id}", handler.HandleFilePatch).Methods("PATCH")
	fileRouter.HandleFunc("/{id}/download", handler.HandleFileDownload).Methods("GET")
	fileRouter.HandleFunc("/{id}/preview", handler.HandleFilePreview).Methods("GET")
	fileRouter.HandleFunc("/", handler.HandleFiles).Methods("GET")
	fileRouter.HandleFunc("/", handler.HandleFileUpload).Methods("POST")
	fileRouter.HandleFunc("/folders/{id}", handler.HandleFolders).Methods("GET")
	fileRouter.HandleFunc("/shared-drives", handler.HandleSharedDrives).Methods("GET")

	fileRouter.Use(middleware.AuthMiddleware(conn))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middleware.CORS(router),
	}

	slog.Info(fmt.Sprintf("starting server on http://localhost:%d", port))

	// Start server
	err = server.ListenAndServe()
	if err != nil {
		slog.Error("error starting server", "err", err)
	}
}
