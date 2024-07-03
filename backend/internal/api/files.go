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

func (s *Config) Files(ctx context.Context, r *http.Request) ([]byte, error) {
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

	switch {
	case err != nil:
		return nil, ErrInternal
	case len(files) == 0:
		return json.Marshal(FilesResponse{
			Data: make([]db.File, 0),
		})
	default:
		return json.Marshal(FilesResponse{
			Data: files,
		})
	}
}

type FoldersResponse struct {
	Data []db.File `json:"data"`
}

func (s *Config) Folders(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	id := r.PathValue("id")

	fileFindByParentIDParams := db.FileFindByParentIDParams{
		ParentID:       pgtype.Text{String: id, Valid: true},
		OrganisationID: user.OrganisationID,
	}

	folder, err := s.DB.FileFindByParentID(ctx, fileFindByParentIDParams)
	if err != nil {
		return nil, ErrInternal
	}

	var resp FoldersResponse
	resp.Data = folder

	if len(folder) == 0 {
		resp.Data = make([]db.File, 0)
	}

	return json.Marshal(resp)
}

type SharedDrivesResponse struct {
	Data []db.File `json:"data"`
}

func (s *Config) SharedDrives(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	drives, err := s.DB.FileFindSharedDrives(ctx, user.OrganisationID)
	if err != nil {
		return nil, ErrInternal
	}

	var resp SharedDrivesResponse
	resp.Data = drives

	if len(drives) == 0 {
		resp.Data = make([]db.File, 0)
	}

	return json.Marshal(resp)
}

type FileUploadResponse struct {
	Data db.File `json:"data"`
}

func (s *Config) FileUpload(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	// TODO: handle the case where file doesn't get uploaded but is already in db
	//  Hint: you can do that with a transaction

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
		return json.Marshal(resp)
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

	fileCreateParams := db.FileCreateParams{
		Name:           header.Filename,
		MimeType:       header.Header.Get("Content-Type"),
		FileSize:       header.Size,
		OrganisationID: user.OrganisationID,
	}

	if r.FormValue("parent_id") != "" {
		fileCreateParams.ParentID = pgtype.Text{String: r.FormValue("parent_id"), Valid: true}
	}

	fileCreated, err := s.DB.FileCreate(ctx, fileCreateParams)
	if err != nil {
		return nil, ErrInternal
	}

	// upload to minio
	_, err = s.MinIO.PutObject(ctx, os.Getenv("MINIO_BUCKET"), fileCreated.ID, file, fileCreateParams.FileSize, minio.PutObjectOptions{})
	if err != nil {
		return nil, ErrInternal
	}

	return json.Marshal(FileUploadResponse{
		Data: fileCreated,
	})
}

func (s *Config) FileDelete(ctx context.Context, r *http.Request) ([]byte, error) {
	_, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	err := s.DB.FileSoftDelete(ctx, r.PathValue("id"))
	if err != nil {
		return nil, ErrInternal
	}

	return nil, nil
}

func (s *Config) FileDownload(w http.ResponseWriter, r *http.Request) {
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

func (s *Config) FilePatch(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	var file db.File
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		return nil, ErrBadRequest
	}

	fileUpdateNameParams := db.FileUpdateNameParams{
		ID:             r.PathValue("id"),
		OrganisationID: user.OrganisationID,
		Name:           file.Name,
	}

	file, err = s.DB.FileUpdateName(ctx, fileUpdateNameParams)
	if err != nil {
		return nil, ErrInternal
	}

	return json.Marshal(file)
}

type FilePreviewResponse struct {
	URL string `json:"url"`
}

func (s *Config) FilePreview(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	id := r.PathValue("id")

	fileFindParams := db.FileFindByIDParams{
		ID:             id,
		OrganisationID: user.OrganisationID,
	}

	file, err := s.DB.FileFindByID(ctx, fileFindParams)
	if err != nil {
		return nil, ErrInternal
	}

	// Get preview url from minio
	presignedURL, err := s.MinIO.PresignedGetObject(ctx, os.Getenv("MINIO_BUCKET"), file.ID, time.Second*60, nil)
	if err != nil {
		return nil, ErrInternal
	}

	return json.Marshal(FilePreviewResponse{
		URL: presignedURL.String(),
	})
}
