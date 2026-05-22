package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rcsabucido/wakealert_api/internal/db"
	"github.com/rcsabucido/wakealert_api/internal/db/sqlc"
)

type medicalHistoryEntryResponse struct {
	HistoryEntryID int64  `json:"history_entry_id"`
	MedicalInfoID  int64  `json:"medical_info_id"`
	Diagnosis      string `json:"diagnosis"`
	IsRecent       bool   `json:"is_recent"`
}

type upsertMedicalHistoryRequest struct {
	MedicalInfoID   int64    `json:"medical_info_id"`
	Diagnoses       []string `json:"diagnoses"`
	MostRecentIndex int      `json:"most_recent_index"`
}

func ListMedicalHistoryByMedicalInfo(w http.ResponseWriter, r *http.Request) {
	medicalInfoIDStr := r.PathValue("medical_info_id")
	medicalInfoID, err := strconv.ParseInt(medicalInfoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid medical info id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	entries, err := q.ListMedicalHistoryByMedicalInfo(r.Context(), medicalInfoID)
	if err != nil {
		http.Error(w, "failed to fetch medical history", http.StatusInternalServerError)
		return
	}

	response := make([]medicalHistoryEntryResponse, 0, len(entries))
	for _, entry := range entries {
		response = append(response, medicalHistoryEntryResponse{
			HistoryEntryID: entry.HistoryEntryID,
			MedicalInfoID:  entry.MedicalInfoID,
			Diagnosis:      entry.Diagnosis,
			IsRecent:       entry.IsRecent.Bool,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpsertMedicalHistoryByMedicalInfo(w http.ResponseWriter, r *http.Request) {
	var req upsertMedicalHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.MedicalInfoID == 0 {
		http.Error(w, "medical_info_id is required", http.StatusBadRequest)
		return
	}
	if len(req.Diagnoses) == 0 {
		http.Error(w, "diagnoses must not be empty", http.StatusBadRequest)
		return
	}
	
	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	response := make([]medicalHistoryEntryResponse, 0, len(req.Diagnoses))
	for idx, diagnosis := range req.Diagnoses {
		if diagnosis == "" {
			http.Error(w, "diagnoses must not contain empty values", http.StatusBadRequest)
			return
		}

		isRecent := idx == req.MostRecentIndex
		existing, err := q.GetMedicalHistoryByInfoAndDiagnosis(r.Context(), sqlc.GetMedicalHistoryByInfoAndDiagnosisParams{
			MedicalInfoID: req.MedicalInfoID,
			Diagnosis:     diagnosis,
		})
		if err != nil {
			if err != pgx.ErrNoRows {
				http.Error(w, "failed to fetch medical history", http.StatusInternalServerError)
				return
			}

			created, err := q.CreateMedicalHistoryEntry(r.Context(), sqlc.CreateMedicalHistoryEntryParams{
				MedicalInfoID: req.MedicalInfoID,
				Diagnosis:     diagnosis,
				IsRecent:      pgtype.Bool{Bool: isRecent, Valid: true},
			})
			if err != nil {
				http.Error(w, "failed to create medical history entry", http.StatusInternalServerError)
				return
			}

			response = append(response, medicalHistoryEntryResponse{
				HistoryEntryID: created.HistoryEntryID,
				MedicalInfoID:  created.MedicalInfoID,
				Diagnosis:      created.Diagnosis,
				IsRecent:       created.IsRecent.Bool,
			})
			continue
		}

		updated, err := q.UpdateMedicalHistoryEntryStatus(r.Context(), sqlc.UpdateMedicalHistoryEntryStatusParams{
			HistoryEntryID: existing.HistoryEntryID,
			IsRecent:       pgtype.Bool{Bool: isRecent, Valid: true},
		})
		if err != nil {
			http.Error(w, "failed to update medical history entry", http.StatusInternalServerError)
			return
		}

		response = append(response, medicalHistoryEntryResponse{
			HistoryEntryID: updated.HistoryEntryID,
			MedicalInfoID:  updated.MedicalInfoID,
			Diagnosis:      updated.Diagnosis,
			IsRecent:       updated.IsRecent.Bool,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SoftDeleteMedicalHistoryByMedicalInfo(w http.ResponseWriter, r *http.Request) {
	medicalInfoIDStr := r.PathValue("medical_info_id")
	medicalInfoID, err := strconv.ParseInt(medicalInfoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid medical info id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	if err := q.SoftDeleteMedicalHistoryByMedicalInfo(r.Context(), medicalInfoID); err != nil {
		http.Error(w, "failed to delete medical history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
