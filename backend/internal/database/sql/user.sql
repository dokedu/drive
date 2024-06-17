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


-- name: UpdateUserConfirmationToken :one
UPDATE users
SET recovery_token = $1
AND recovery_sent_at = $2
WHERE id = $3
  AND deleted_at IS NULL
RETURNING *;

-- name: ResetUserConfirmationToken :one
UPDATE users
SET recovery_token = NULL
AND recovery_sent_at = NULL
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
