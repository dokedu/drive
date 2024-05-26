package api

import (
	"encoding/json"
	"example/internal/database/db"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
)

type FilesResponse struct {
	Data []db.File `json:"data"`
}

func (s *Config) HandleFiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parentId := r.URL.Query().Get("parent_id")
	sharedDrives := r.URL.Query().Get("shared_drive")

	var files []db.File
	var err error
	if sharedDrives != "" {
		files, err = s.DB.FileFindSharedDrives(ctx, "1")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if parentId == "" {
		files, err = s.DB.FileFindAll(ctx, "1")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		files, err = s.DB.FileFindByParentID(ctx, db.FileFindByParentIDParams{
			ParentID:       pgtype.Text{String: parentId, Valid: true},
			OrganisationID: "1",
		})
	}

	var resp FilesResponse
	resp.Data = files

	if len(files) == 0 {
		resp.Data = make([]db.File, 0)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Config) HandleFolders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	folder, err := s.DB.FileFindByParentID(ctx, db.FileFindByParentIDParams{
		ParentID:       pgtype.Text{String: id, Valid: true},
		OrganisationID: "1",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp FilesResponse
	resp.Data = folder

	if len(folder) == 0 {
		resp.Data = make([]db.File, 0)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Config) HandleSharedDrives(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	drives, err := s.DB.FileFindSharedDrives(ctx, "1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp FilesResponse
	resp.Data = drives

	if len(drives) == 0 {
		resp.Data = make([]db.File, 0)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

type FileUploadResponse struct {
	Data db.File `json:"data"`
}

// TODO: handle the case where file doesn't get uploaded but is already in db
func (s *Config) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.FormValue("is_folder") != "" {
		folder, err := s.DB.FileCreateFolder(ctx, db.FileCreateFolderParams{
			Name:           r.FormValue("name"),
			OrganisationID: "1",
		})
		if err != nil {
			http.Error(w, "Error creating file record: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := FileUploadResponse{
			Data: folder,
		}
		responseBody, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Error marshaling response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(responseBody)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileName := header.Filename
	mimeType := header.Header.Get("Content-Type")
	fileSize := header.Size

	fileCreateParams := db.FileCreateParams{
		Name:           fileName,
		MimeType:       mimeType,
		FileSize:       fileSize,
		OrganisationID: "1",
	}

	if r.FormValue("parent_id") != "" {
		fileCreateParams.ParentID = pgtype.Text{String: r.FormValue("parent_id"), Valid: true}
	}

	fileCreated, err := s.DB.FileCreate(ctx, fileCreateParams)
	if err != nil {
		http.Error(w, "Error creating file record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := FileUploadResponse{
		Data: fileCreated,
	}
	responseBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error marshaling response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// upload to minio
	_, err = s.MinIO.PutObject(ctx, "drive", fileCreated.ID, file, fileSize, minio.PutObjectOptions{})
	if err != nil {
		http.Error(w, "Error uploading file to minio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseBody)
}

func (s *Config) HandleFileDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")

	err := s.DB.FileSoftDelete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Config) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")

	file, err := s.DB.FileFindByID(ctx, db.FileFindByIDParams{
		ID:             id,
		OrganisationID: "1",
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	// download from minio
	object, err := s.MinIO.GetObject(ctx, "drive", file.ID, minio.GetObjectOptions{})
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

func (s *Config) HandleFilePatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")

	var file db.File

	// parse json body
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("the file", "file", file)

	file, err = s.DB.FileUpdateName(ctx, db.FileUpdateNameParams{
		ID:             id,
		OrganisationID: "1",
		Name:           file.Name,
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(file)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Config) HandleFilePreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	// Get file from db
	file, err := s.DB.FileFindByID(ctx, db.FileFindByIDParams{
		ID:             id,
		OrganisationID: "1",
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get preview url from minio
	presignedURL, err := s.MinIO.PresignedGetObject(ctx, "drive", file.ID, time.Second*60, nil)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct {
		URL string `json:"url"`
	}{}
	res.URL = presignedURL.String()
	body, err := json.Marshal(res)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
