-- name: GetAlert :one
SELECT * FROM alert
WHERE alert_id = $1 AND is_deleted = FALSE;

-- name: ListAlerts :many
SELECT * FROM alert
where is_deleted = FALSE
ORDER BY alert_time DESC;

-- name: ListAlertsByVictim :many
SELECT * FROM alert
WHERE victim_id = $1 AND is_deleted = FALSE
ORDER BY alert_time DESC;

-- name: CreateAlert :one
INSERT INTO alert (latitude, longitude, victim_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CompleteAlert :one
UPDATE alert
SET is_completed = TRUE
WHERE alert_id = $1
RETURNING *;

-- name: SoftDeleteAlert :one
UPDATE alert
SET is_deleted = TRUE
WHERE alert_id = $1
RETURNING *;

-- name: UpdateAlertStatus :one
UPDATE alert
SET is_completed = $2
WHERE alert_id = $1
RETURNING *;