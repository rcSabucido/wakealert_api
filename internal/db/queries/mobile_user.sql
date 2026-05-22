-- name: GetMobileUser :one
SELECT * FROM mobile_user
WHERE mobile_user_id = $1;

-- name: GetMobileUserByEmail :one
SELECT * FROM mobile_user
WHERE email = $1;

-- name: CreateMobileUser :one
INSERT INTO mobile_user (email, password)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateMobileUserPassword :one
UPDATE mobile_user
SET password = $1
WHERE mobile_user_id = $2
RETURNING *;

-- name: DeleteMobileUser :exec
DELETE FROM mobile_user
WHERE mobile_user_id = $1;