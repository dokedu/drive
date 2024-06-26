-- name: UserFind :one
SELECT *
FROM users
WHERE id = $1
  AND organisation_id = $2
  AND deleted_at IS NULL;

-- name: UserFindByEmail :one
SELECT *
FROM users
WHERE email = $1
  AND deleted_at IS NULL;

-- name: UserFindByID :one
SELECT *
FROM users
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateUserConfirmationToken :one
UPDATE users
SET recovery_token = @recovery_token, recovery_sent_at = @recovery_sent_at
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: ResetUserConfirmationToken :one
UPDATE users
SET recovery_token = NULL, recovery_sent_at = NULL
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: CreateOrganisation :one
INSERT INTO organisations (name) VALUES ($1) RETURNING *;

-- name: OrganisationFindByName :one
SELECT *
FROM organisations
WHERE name = $1
  AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name, organisation_id, role) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UserFindByToken :one
SELECT *
FROM users
WHERE recovery_token = $1
  AND deleted_at IS NULL;

-- name: CreateSession :one
INSERT INTO sessions (user_id, token ) VALUES ($1, $2) RETURNING *;

-- name: FindSessionByToken :one
SELECT *
FROM sessions
WHERE token = $1
  AND deleted_at IS NULL;

-- name: RemoveSession :one
UPDATE sessions
SET deleted_at = NOW()
WHERE token = $1
  AND deleted_at IS NULL
RETURNING *;
