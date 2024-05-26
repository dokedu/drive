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