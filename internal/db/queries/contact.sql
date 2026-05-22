-- name: GetContact :one
SELECT * FROM contact
WHERE contact_id = $1;

-- name: ListContactsByMobileUser :many
SELECT * FROM contact
WHERE mobile_user_id = $1;

-- name: GetPrimaryContact :one
SELECT * FROM contact
WHERE mobile_user_id = $1 AND is_primary = TRUE;

-- name: CreateContact :one
INSERT INTO contact (
    mobile_user_id,
    first_name,
    last_name,
    phone_number,
    relationship_type_id,
    is_primary
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateContact :one
UPDATE contact
SET
    first_name = $1,
    last_name = $2,
    phone_number = $3,
    relationship_type_id = $4,
    is_primary = $5
WHERE contact_id = $6
RETURNING *;

-- name: ClearPrimaryContactsByMobileUser :exec
UPDATE contact
SET is_primary = FALSE
WHERE mobile_user_id = $1;

-- name: DeleteContact :exec
UPDATE contact
SET is_deleted = TRUE
WHERE mobile_user_id = $1 AND contact_id = $2;

-- name: DeleteContactByDetails :exec
UPDATE contact AS c
SET is_deleted = TRUE
FROM relationship_type AS rt
WHERE c.mobile_user_id = $1
    AND c.is_deleted = FALSE
    AND c.first_name = $2
    AND c.last_name = $3
    AND c.phone_number = $4
    AND rt.relationship_name = $5
    AND c.relationship_type_id = rt.relationship_type_id;

-- name: UpdateContactByDetails :one
UPDATE contact AS c
SET
    first_name = $7,
    last_name = $8,
    phone_number = $9,
    relationship_type_id = rt_new.relationship_type_id,
    is_primary = $11
FROM relationship_type AS rt_old,
     relationship_type AS rt_new
WHERE c.mobile_user_id = $1
    AND c.is_deleted = FALSE
    AND c.first_name = $2
    AND c.last_name = $3
    AND c.phone_number = $4
    AND c.is_primary = $6
    AND rt_old.relationship_name = $5
    AND c.relationship_type_id = rt_old.relationship_type_id
    AND rt_new.relationship_name = $10
RETURNING *;

-- name: GetRelationshipType :one
SELECT * FROM relationship_type
WHERE relationship_type_id = $1;