-- name: GetMedicalInformation :one
SELECT * FROM medical_information
WHERE medical_info_id = $1;

-- name: CreateMedicalInformation :one
INSERT INTO medical_information (
    allergies,
    medication,
    medical_notes,
    pregnancy_status_id,
    donor_status_id,
    blood_type_id,
    last_diagnosis_date,
    last_diagnosis_hospital_name
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateMedicalInformation :one
UPDATE medical_information
SET
    allergies = $1,
    medication = $2,
    medical_notes = $3,
    pregnancy_status_id = $4,
    donor_status_id = $5,
    blood_type_id = $6,
    last_diagnosis_date = $7,
    last_diagnosis_hospital_name = $8
WHERE medical_info_id = $9
RETURNING *;

-- name: UpdateMedicalInformationByNames :one
UPDATE medical_information AS mi
SET
    allergies = $2,
    medication = $3,
    medical_notes = $4,
    pregnancy_status_id = ps.pregnancy_status_id,
    donor_status_id = ods.donor_status_id,
    blood_type_id = bt.blood_type_id,
    last_diagnosis_date = $8,
    last_diagnosis_hospital_name = $9
FROM pregnancy_status AS ps,
     organ_donor_status AS ods,
     blood_type AS bt
WHERE mi.medical_info_id = $1
    AND ps.pregnancy_status_name = $5
    AND ods.donor_status_name = $6
    AND bt.blood_type_name = $7
RETURNING *;

-- name: GetPregnancyStatus :one
SELECT * FROM pregnancy_status
WHERE pregnancy_status_id = $1;

-- name: GetOrganDonorStatus :one
SELECT * FROM organ_donor_status
WHERE donor_status_id = $1;

-- name: GetBloodType :one
SELECT * FROM blood_type
WHERE blood_type_id = $1;