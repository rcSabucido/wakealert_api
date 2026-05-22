-- name: GetVictim :one
SELECT * FROM victim
WHERE victim_id = $1;

-- name: GetVictimByMobileUser :one
SELECT * FROM victim
WHERE mobile_user_id = $1;

-- name: GetVictimAddressID :one
SELECT address_id FROM victim
WHERE victim_id = $1;

-- name: UpdateVictimAddressID :one
UPDATE victim
SET address_id = $1
WHERE victim_id = $2
RETURNING address_id;

-- name: ListVictims :many
SELECT * FROM victim
ORDER BY victim_id;

-- name: CreateVictim :one
INSERT INTO victim (
    mobile_user_id,
    first_name,
    last_name,
    -- birth_date,
    -- address_id,
    medical_info_id
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateVictim :one
UPDATE victim
SET
    first_name = $1,
    last_name = $2,
    birth_date = $3,
    address_id = $4
WHERE victim_id = $5
RETURNING *;

-- name: DeleteVictim :exec
DELETE FROM victim
WHERE victim_id = $1;

-- name: GetVictimDetails :one
SELECT
    v.victim_id,
    v.first_name,
    v.last_name,
    v.birth_date,
    v.mobile_user_id,
    mi.allergies,
    mi.medication,
    mi.medical_notes,
    mi.last_diagnosis_date,
    mi.last_diagnosis_hospital_name,
    ps.pregnancy_status_name,
    bt.blood_type_name,
    ods.donor_status_name,
    COALESCE(concat_ws(
        ', ',
        nullif(al.address_line, ''),
        nullif(b.barangay_name, ''),
        nullif(moc.city_mun_name, ''),
        nullif(poh.province_or_huc_name, ''),
        nullif(r.region_name, '')
    ), '') AS address_text,
    pc.phone_number AS primary_contact_phone,
    rt.relationship_name AS primary_contact_relationship,
    COALESCE((
        SELECT mhe2.diagnosis
        FROM medical_history_entry mhe2
        WHERE mhe2.medical_info_id = mi.medical_info_id
        ORDER BY coalesce(mhe2.is_recent, false) DESC, mhe2.history_entry_id ASC
        LIMIT 1
    ), '') AS last_diagnosis,
    COALESCE(
        string_agg(mhe.diagnosis, ', ' ORDER BY mhe.history_entry_id),
        ''
    ) AS medical_history
FROM victim v
JOIN medical_information mi ON mi.medical_info_id = v.medical_info_id
LEFT JOIN blood_type bt ON mi.blood_type_id = bt.blood_type_id
LEFT JOIN pregnancy_status ps ON ps.pregnancy_status_id = mi.pregnancy_status_id
LEFT JOIN organ_donor_status ods ON ods.donor_status_id = mi.donor_status_id
LEFT JOIN address_line al ON al.address_id = v.address_id

LEFT JOIN barangays b ON b.barangay_psgc = al.barangay_psgc
LEFT JOIN municipal_or_city moc ON moc.city_mun_psgc = b.city_mun_psgc
LEFT JOIN province_or_huc poh ON poh.province_or_huc_psgc = moc.province_or_huc_psgc
LEFT JOIN region r ON r.region_psgc = poh.region_psgc

LEFT JOIN contact pc ON pc.mobile_user_id = v.mobile_user_id AND pc.is_primary = TRUE
LEFT JOIN relationship_type rt ON rt.relationship_type_id = pc.relationship_type_id
LEFT JOIN medical_history_entry mhe
  ON mhe.medical_info_id = mi.medical_info_id
 AND mhe.is_deleted = FALSE
WHERE v.victim_id = $1 --AND mhe.is_deleted = FALSE
GROUP BY
    v.victim_id,
    v.first_name,
    v.last_name,
    v.birth_date,
    v.mobile_user_id,
    mi.medical_info_id,
    mi.allergies,
    mi.medication,
    mi.medical_notes,
    mi.last_diagnosis_date,
    mi.last_diagnosis_hospital_name,
    ps.pregnancy_status_name,
    bt.blood_type_name,
    ods.donor_status_name,
    al.address_line,
    b.barangay_name,
    moc.city_mun_name,
    poh.province_or_huc_name,
    r.region_name,
    pc.phone_number,
    rt.relationship_name;

