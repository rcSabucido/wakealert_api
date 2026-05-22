-- name: GetReceiverUser :one
SELECT * FROM receiver_user
WHERE receiver_user_id = $1;

-- name: GetReceiverUserByUsername :one
SELECT * FROM receiver_user
WHERE username = $1;

-- name: CreateReceiverUser :one
INSERT INTO receiver_user (username, password)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteReceiverUser :exec
DELETE FROM receiver_user
WHERE receiver_user_id = $1;