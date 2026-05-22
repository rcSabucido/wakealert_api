package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rcsabucido/wakealert_api/internal/db"
	"github.com/rcsabucido/wakealert_api/internal/db/sqlc"
)

type updateMedicalInformationRequest struct {
	Allergies                 *string `json:"allergies"`
	Medication                *string `json:"medication"`
	MedicalNotes              *string `json:"medical_notes"`
	LastDiagnosisDate         *string `json:"last_diagnosis_date"`
	LastDiagnosisHospitalName *string `json:"last_diagnosis_hospital_name"`
	PregnancyStatusName       *string `json:"pregnancy_status_name"`
	DonorStatusName           *string `json:"donor_status_name"`
	BloodTypeName             *string `json:"blood_type_name"`
}

func CreateMedicalInformation(w http.ResponseWriter, r *http.Request) {
	// No request body needed; create with default lookup values and null nullable fields.
	pool, err := db.Pool(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	var pregnancyStatusID int64
	if err := pool.QueryRow(r.Context(),
		"SELECT pregnancy_status_id FROM pregnancy_status WHERE pregnancy_status_name = $1",
		"Unknown",
	).Scan(&pregnancyStatusID); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "default pregnancy status not found", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to resolve pregnancy status", http.StatusInternalServerError)
		return
	}

	var donorStatusID int64
	if err := pool.QueryRow(r.Context(),
		"SELECT donor_status_id FROM organ_donor_status WHERE donor_status_name = $1",
		"Unknown",
	).Scan(&donorStatusID); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "default donor status not found", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to resolve donor status", http.StatusInternalServerError)
		return
	}

	var bloodTypeID int64
	if err := pool.QueryRow(r.Context(),
		"SELECT blood_type_id FROM blood_type WHERE blood_type_name = $1",
		"Unknown",
	).Scan(&bloodTypeID); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "default blood type not found", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to resolve blood type", http.StatusInternalServerError)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	info, err := q.CreateMedicalInformation(r.Context(), sqlc.CreateMedicalInformationParams{
		Allergies:                 pgtype.Text{Valid: false},
		Medication:                pgtype.Text{Valid: false},
		MedicalNotes:              pgtype.Text{Valid: false},
		PregnancyStatusID:         pregnancyStatusID,
		DonorStatusID:             donorStatusID,
		BloodTypeID:               bloodTypeID,
		LastDiagnosisDate:         pgtype.Date{Valid: false},
		LastDiagnosisHospitalName: pgtype.Text{Valid: false},
	})
	if err != nil {
		http.Error(w, "failed to create medical information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		MedicalInfoID int64 `json:"medical_info_id"`
	}{
		MedicalInfoID: info.MedicalInfoID,
	})
}

func UpdateMedicalInformation(w http.ResponseWriter, r *http.Request) {
	medicalInfoIDStr := r.PathValue("medical_info_id")
	medicalInfoID, err := strconv.ParseInt(medicalInfoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid medical info id", http.StatusBadRequest)
		return
	}

	var req updateMedicalInformationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Allergies == nil && req.Medication == nil && req.MedicalNotes == nil && req.LastDiagnosisDate == nil && req.LastDiagnosisHospitalName == nil && req.PregnancyStatusName == nil && req.DonorStatusName == nil && req.BloodTypeName == nil {
		http.Error(w, "at least one field must be provided", http.StatusBadRequest)
		return
	}
	if req.Allergies != nil && *req.Allergies == "" {
		http.Error(w, "allergies cannot be empty", http.StatusBadRequest)
		return
	}
	if req.Medication != nil && *req.Medication == "" {
		http.Error(w, "medication cannot be empty", http.StatusBadRequest)
		return
	}
	if req.MedicalNotes != nil && *req.MedicalNotes == "" {
		http.Error(w, "medical_notes cannot be empty", http.StatusBadRequest)
		return
	}
	if req.LastDiagnosisDate != nil && *req.LastDiagnosisDate == "" {
		http.Error(w, "last_diagnosis_date cannot be empty", http.StatusBadRequest)
		return
	}
	if req.LastDiagnosisHospitalName != nil && *req.LastDiagnosisHospitalName == "" {
		http.Error(w, "last_diagnosis_hospital_name cannot be empty", http.StatusBadRequest)
		return
	}
	if req.PregnancyStatusName != nil && *req.PregnancyStatusName == "" {
		http.Error(w, "pregnancy_status_name cannot be empty", http.StatusBadRequest)
		return
	}
	if req.DonorStatusName != nil && *req.DonorStatusName == "" {
		http.Error(w, "donor_status_name cannot be empty", http.StatusBadRequest)
		return
	}
	if req.BloodTypeName != nil && *req.BloodTypeName == "" {
		http.Error(w, "blood_type_name cannot be empty", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	existing, err := q.GetMedicalInformation(r.Context(), medicalInfoID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "medical information not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch medical information", http.StatusInternalServerError)
		return
	}

	allergies := existing.Allergies
	if req.Allergies != nil {
		allergies = pgtype.Text{String: *req.Allergies, Valid: true}
	}
	medication := existing.Medication
	if req.Medication != nil {
		medication = pgtype.Text{String: *req.Medication, Valid: true}
	}
	medicalNotes := existing.MedicalNotes
	if req.MedicalNotes != nil {
		medicalNotes = pgtype.Text{String: *req.MedicalNotes, Valid: true}
	}
	lastDiagnosisDate := existing.LastDiagnosisDate
	if req.LastDiagnosisDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.LastDiagnosisDate)
		if err != nil {
			http.Error(w, "last_diagnosis_date must be YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		lastDiagnosisDate = pgtype.Date{Time: parsed, Valid: true}
	}
	lastDiagnosisHospitalName := existing.LastDiagnosisHospitalName
	if req.LastDiagnosisHospitalName != nil {
		lastDiagnosisHospitalName = pgtype.Text{String: *req.LastDiagnosisHospitalName, Valid: true}
	}

	pregnancyStatusName := ""
	if req.PregnancyStatusName != nil {
		pregnancyStatusName = *req.PregnancyStatusName
	} else {
		status, err := q.GetPregnancyStatus(r.Context(), existing.PregnancyStatusID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "pregnancy status not found", http.StatusInternalServerError)
				return
			}
			http.Error(w, "failed to resolve pregnancy status", http.StatusInternalServerError)
			return
		}
		pregnancyStatusName = status.PregnancyStatusName
	}

	donorStatusName := ""
	if req.DonorStatusName != nil {
		donorStatusName = *req.DonorStatusName
	} else {
		status, err := q.GetOrganDonorStatus(r.Context(), existing.DonorStatusID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "donor status not found", http.StatusInternalServerError)
				return
			}
			http.Error(w, "failed to resolve donor status", http.StatusInternalServerError)
			return
		}
		donorStatusName = status.DonorStatusName
	}

	bloodTypeName := ""
	if req.BloodTypeName != nil {
		bloodTypeName = *req.BloodTypeName
	} else {
		bloodType, err := q.GetBloodType(r.Context(), existing.BloodTypeID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "blood type not found", http.StatusInternalServerError)
				return
			}
			http.Error(w, "failed to resolve blood type", http.StatusInternalServerError)
			return
		}
		bloodTypeName = bloodType.BloodTypeName
	}

	_, err = q.UpdateMedicalInformationByNames(r.Context(), sqlc.UpdateMedicalInformationByNamesParams{
		MedicalInfoID:             medicalInfoID,
		Allergies:                 allergies,
		Medication:                medication,
		MedicalNotes:              medicalNotes,
		PregnancyStatusName:       pregnancyStatusName,
		DonorStatusName:           donorStatusName,
		BloodTypeName:             bloodTypeName,
		LastDiagnosisDate:         lastDiagnosisDate,
		LastDiagnosisHospitalName: lastDiagnosisHospitalName,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "invalid pregnancy_status_name, donor_status_name, or blood_type_name", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to update medical information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		MedicalInfoID int64 `json:"medical_info_id"`
	}{
		MedicalInfoID: medicalInfoID,
	})
}
