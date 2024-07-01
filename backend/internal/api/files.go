package api

import (
	"context"
	"encoding/json"
	"example/internal/database/db"
	"example/internal/middleware"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
)

type FilesResponse struct {
	Data []db.File `json:"data"`
}

func (s *Config) HandleFiles(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	parentId := r.URL.Query().Get("parent_id")
	sharedDrives := r.URL.Query().Get("shared_drive")

	var files []db.File
	var err error
	if sharedDrives != "" {
		files, err = s.DB.FileFindSharedDrives(ctx, user.OrganisationID)
		if err != nil {
			return nil, ErrInternal
		}
	} else if parentId == "" {
		files, err = s.DB.FileFindAll(ctx, user.OrganisationID)
		if err != nil {
			return nil, ErrInternal
		}
	} else {
		files, err = s.DB.FileFindByParentID(ctx, db.FileFindByParentIDParams{
			ParentID:       pgtype.Text{String: parentId, Valid: true},
			OrganisationID: user.OrganisationID,
		})
	}

	var resp FilesResponse
	resp.Data = files

	if len(files) == 0 {
		resp.Data = make([]db.File, 0)
	}

	return json.Marshal(resp)
}

func (s *Config) HandleFolders(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	id := r.PathValue("id")

	folder, err := s.DB.FileFindByParentID(ctx, db.FileFindByParentIDParams{
		ParentID:       pgtype.Text{String: id, Valid: true},
		OrganisationID: user.OrganisationID,
	})
	if err != nil {
		return nil, ErrInternal
	}

	var resp FilesResponse
	resp.Data = folder

	if len(folder) == 0 {
		resp.Data = make([]db.File, 0)
	}

	return json.Marshal(resp)
}

func (s *Config) HandleSharedDrives(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	drives, err := s.DB.FileFindSharedDrives(ctx, user.OrganisationID)
	if err != nil {
		return nil, ErrInternal
	}

	var resp FilesResponse
	resp.Data = drives

	if len(drives) == 0 {
		resp.Data = make([]db.File, 0)
	}

	return json.Marshal(resp)
}

type FileUploadResponse struct {
	Data db.File `json:"data"`
}

// TODO: handle the case where file doesn't get uploaded but is already in db
func (s *Config) HandleFileUpload(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	if r.FormValue("is_folder") != "" {
		folder, err := s.DB.FileCreateFolder(ctx, db.FileCreateFolderParams{
			Name:           r.FormValue("name"),
			OrganisationID: user.OrganisationID,
		})
		if err != nil {
			return nil, ErrInternal
		}

		resp := FileUploadResponse{
			Data: folder,
		}
		responseBody, err := json.Marshal(resp)
		if err != nil {
			return nil, ErrInternal
		}

		return responseBody, nil
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, ErrInternal
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, ErrInternal
	}
	defer file.Close()

	fileName := header.Filename
	mimeType := header.Header.Get("Content-Type")
	fileSize := header.Size

	fileCreateParams := db.FileCreateParams{
		Name:           fileName,
		MimeType:       mimeType,
		FileSize:       fileSize,
		OrganisationID: user.OrganisationID,
	}

	if r.FormValue("parent_id") != "" {
		fileCreateParams.ParentID = pgtype.Text{String: r.FormValue("parent_id"), Valid: true}
	}

	fileCreated, err := s.DB.FileCreate(ctx, fileCreateParams)
	if err != nil {
		return nil, ErrInternal
	}

	resp := FileUploadResponse{
		Data: fileCreated,
	}
	responseBody, err := json.Marshal(resp)
	if err != nil {
		return nil, ErrInternal
	}

	// upload to minio
	_, err = s.MinIO.PutObject(ctx, os.Getenv("MINIO_BUCKET"), fileCreated.ID, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		return nil, ErrInternal
	}

	return responseBody, nil
}

func (s *Config) HandleFileDelete(ctx context.Context, r *http.Request) ([]byte, error) {
	_, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	id := r.PathValue("id")

	err := s.DB.FileSoftDelete(ctx, id)
	if err != nil {
		return nil, ErrInternal
	}

	return nil, nil
}

func (s *Config) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	file, err := s.DB.FileFindByID(ctx, db.FileFindByIDParams{
		ID:             id,
		OrganisationID: user.OrganisationID,
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	// download from minio
	object, err := s.MinIO.GetObject(ctx, os.Getenv("MINIO_BUCKET"), file.ID, minio.GetObjectOptions{})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer object.Close()

	_, err = io.Copy(w, object)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func (s *Config) HandleFilePatch(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	id := r.PathValue("id")

	var file db.File

	// parse json body
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		slog.Error(err.Error())
		return nil, ErrBadRequest
	}

	slog.Info("the file", "file", file)

	file, err = s.DB.FileUpdateName(ctx, db.FileUpdateNameParams{
		ID:             id,
		OrganisationID: user.OrganisationID,
		Name:           file.Name,
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, ErrInternal
	}

	return json.Marshal(file)
}

func (s *Config) HandleFilePreview(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}
	id := r.PathValue("id")

	// Get file from db
	file, err := s.DB.FileFindByID(ctx, db.FileFindByIDParams{
		ID:             id,
		OrganisationID: user.OrganisationID,
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, ErrInternal
	}

	// Get preview url from minio
	presignedURL, err := s.MinIO.PresignedGetObject(ctx, os.Getenv("MINIO_BUCKET"), file.ID, time.Second*60, nil)
	if err != nil {
		slog.Error(err.Error())
		return nil, ErrInternal
	}

	res := struct {
		URL string `json:"url"`
	}{}
	res.URL = presignedURL.String()

	return json.Marshal(res)
}
