-- name: FileFindAll :many
SELECT *
FROM files
WHERE deleted_at IS NULL
  AND parent_id IS NULL
  AND shared_drive IS FALSE
  AND organisation_id = $1
  AND deleted_at IS NULL
ORDER BY is_folder DESC, name;

-- name: FileFindByParentID :many
SELECT *
FROM files
WHERE parent_id = $1
  AND organisation_id = $2
  AND deleted_at IS NULL
ORDER BY is_folder, name;

-- name: FileFindSharedDrives :many
SELECT *
FROM files
WHERE shared_drive IS TRUE
  AND organisation_id = $1
  AND deleted_at IS NULL
ORDER BY is_folder, name;

-- name: FileCreate :one
INSERT INTO files (name, mime_type, file_size, parent_id, organisation_id)
VALUES (@name, @mime_type, @file_size, @parent_id, @organisation_id)
RETURNING *;

-- name: FileCreateFolder :one
INSERT INTO files (name, mime_type, file_size, is_folder, organisation_id)
VALUES (@name, 'directory', 0, TRUE, @organisation_id)
RETURNING *;

-- name: FileFindTrashed :many
SELECT *
FROM files
WHERE deleted_at IS NOT NULL
  AND organisation_id = $1;

-- name: FileSoftDelete :exec
UPDATE files
SET deleted_at = NOW()
WHERE id = $1;

-- name: FileFindByID :one
SELECT *
FROM files
WHERE id = $1
  AND organisation_id = $2
  AND deleted_at IS NULL;

-- name: FileUpdateName :one
UPDATE files
SET name = $1
WHERE id = $2
  AND organisation_id = $3
  AND deleted_at IS NULL
RETURNING *;