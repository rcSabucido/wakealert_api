-- name: GetMedicalHistoryEntry :one
SELECT * FROM medical_history_entry
WHERE history_entry_id = $1;

-- name: GetMedicalHistoryByInfoAndDiagnosis :one
SELECT * FROM medical_history_entry
WHERE medical_info_id = $1 AND diagnosis = $2;

-- name: ListMedicalHistoryByMedicalInfo :many
SELECT * FROM medical_history_entry
WHERE medical_info_id = $1 AND is_deleted = FALSE;

-- name: CreateMedicalHistoryEntry :one
INSERT INTO medical_history_entry (medical_info_id, diagnosis, is_recent)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateMedicalHistoryByMedicalInfo :many
UPDATE medical_history_entry
SET
	diagnosis = COALESCE($2, diagnosis),
	is_recent = COALESCE($3, is_recent)
WHERE medical_info_id = $1
RETURNING *;

-- name: UpdateMedicalHistoryEntryStatus :one
UPDATE medical_history_entry
SET
    is_deleted = FALSE,
    is_recent = $2
WHERE history_entry_id = $1
RETURNING *;

-- name: SoftDeleteMedicalHistoryByMedicalInfo :exec
UPDATE medical_history_entry
SET is_deleted = TRUE
WHERE medical_info_id = $1;

-- name: DeleteMedicalHistoryEntry :exec
DELETE FROM medical_history_entry
WHERE history_entry_id = $1;
