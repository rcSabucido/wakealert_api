package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rcsabucido/wakealert_api/internal/db"
	"github.com/rcsabucido/wakealert_api/internal/db/sqlc"
)

type createVictimRequest struct {
	MobileUserID  int64  `json:"mobile_user_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MedicalInfoID int64  `json:"medical_info_id"`
}

type updateVictimRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	BirthDate *string `json:"birth_date"`
}

type updateVictimAddressRequest struct {
	AddressID *int64 `json:"address_id"`
}

type victimPartialResponse struct {
	VictimID      int64  `json:"victim_id"`
	MobileUserID  int64  `json:"mobile_user_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MedicalInfoID int64  `json:"medical_info_id"`
}

type victimAddressResponse struct {
	VictimID  int64  `json:"victim_id"`
	AddressID *int64 `json:"address_id"`
}

type VictimDetailsResponse struct {
	VictimID         int64  `json:"victimId"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	FullName         string `json:"fullName"`
	BirthDate        string `json:"birthDate"`
	Age              int    `json:"age"`
	PrimaryContact   string `json:"primaryContact"`
	Address          string `json:"address"`
	PregnancyStatus  string `json:"pregnancyStatus"`
	OrganDonor       string `json:"organDonor"`
	BloodType        string `json:"bloodType"`
	LastDiagnosis    string `json:"lastDiagnosis"`
	DiagnosisDate    string `json:"diagnosisDate"`
	PlaceOfDiagnosis string `json:"placeOfDiagnosis"`
	Allergies        string `json:"allergies"`
	Medication       string `json:"medication"`
	MedicalHistory   string `json:"medicalHistory"`
	MedicalNote      string `json:"medicalNote"`
}

func interfaceString(v interface{}) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

func ListVictims(w http.ResponseWriter, r *http.Request) {
	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	victims, err := q.ListVictims(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch victims", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(victims)
}

func CreateVictim(w http.ResponseWriter, r *http.Request) {
	var req createVictimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.MobileUserID == 0 || req.FirstName == "" || req.LastName == "" || req.MedicalInfoID == 0 {
		http.Error(w, "mobile_user_id, first_name, last_name, birth_date, and medical_info_id are required", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	victim, err := q.CreateVictim(r.Context(), sqlc.CreateVictimParams{
		MobileUserID:  req.MobileUserID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		MedicalInfoID: req.MedicalInfoID,
	})
	if err != nil {
		http.Error(w, "failed to create victim", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(victimPartialResponse{
		VictimID:      victim.VictimID,
		MobileUserID:  victim.MobileUserID,
		FirstName:     victim.FirstName,
		LastName:      victim.LastName,
		MedicalInfoID: victim.MedicalInfoID,
	})
}

func GetVictimByMobileUser(w http.ResponseWriter, r *http.Request) {
	mobileUserIDStr := r.PathValue("mobile_user_id")
	mobileUserID, err := strconv.ParseInt(mobileUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid mobile user id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	victim, err := q.GetVictimByMobileUser(r.Context(), mobileUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "victim not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch victim", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(victimPartialResponse{
		VictimID:      victim.VictimID,
		MobileUserID:  victim.MobileUserID,
		FirstName:     victim.FirstName,
		LastName:      victim.LastName,
		MedicalInfoID: victim.MedicalInfoID,
	})
}

func GetVictimAddressID(w http.ResponseWriter, r *http.Request) {
	victimIDStr := r.PathValue("victim_id")
	victimID, err := strconv.ParseInt(victimIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid victim id", http.StatusBadRequest)
		return
	}

	pool, err := db.Pool(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	var addressID pgtype.Int8
	if err := pool.QueryRow(r.Context(),
		"SELECT address_id FROM victim WHERE victim_id = $1",
		victimID,
	).Scan(&addressID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "victim not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch victim address", http.StatusInternalServerError)
		return
	}

	var addressIDValue *int64
	if addressID.Valid {
		addressIDValue = &addressID.Int64
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(victimAddressResponse{
		VictimID:  victimID,
		AddressID: addressIDValue,
	})
}

func UpdateVictimAddressID(w http.ResponseWriter, r *http.Request) {
	victimIDStr := r.PathValue("victim_id")
	victimID, err := strconv.ParseInt(victimIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid victim id", http.StatusBadRequest)
		return
	}

	var req updateVictimAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	addressID := pgtype.Int8{Valid: false}
	if req.AddressID != nil {
		addressID = pgtype.Int8{Int64: *req.AddressID, Valid: true}
	}

	pool, err := db.Pool(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	if err := pool.QueryRow(r.Context(),
		"UPDATE victim SET address_id = $1 WHERE victim_id = $2 RETURNING address_id",
		addressID,
		victimID,
	).Scan(&addressID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "victim not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update victim address", http.StatusInternalServerError)
		return
	}

	var addressIDValue *int64
	if addressID.Valid {
		addressIDValue = &addressID.Int64
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(victimAddressResponse{
		VictimID:  victimID,
		AddressID: addressIDValue,
	})
}

func UpdateVictim(w http.ResponseWriter, r *http.Request) {
	victimIDStr := r.PathValue("victim_id")
	victimID, err := strconv.ParseInt(victimIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid victim id", http.StatusBadRequest)
		return
	}

	var req updateVictimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.FirstName == nil && req.LastName == nil && req.BirthDate == nil {
		http.Error(w, "at least one field must be provided", http.StatusBadRequest)
		return
	}
	if req.FirstName != nil && *req.FirstName == "" {
		http.Error(w, "first_name cannot be empty", http.StatusBadRequest)
		return
	}
	if req.LastName != nil && *req.LastName == "" {
		http.Error(w, "last_name cannot be empty", http.StatusBadRequest)
		return
	}
	if req.BirthDate != nil && *req.BirthDate == "" {
		http.Error(w, "birth_date cannot be empty", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	existing, err := q.GetVictim(r.Context(), victimID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "victim not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch victim", http.StatusInternalServerError)
		return
	}

	firstName := existing.FirstName
	if req.FirstName != nil {
		firstName = *req.FirstName
	}
	lastName := existing.LastName
	if req.LastName != nil {
		lastName = *req.LastName
	}
	birthDate := existing.BirthDate
	if req.BirthDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			http.Error(w, "birth_date must be YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		birthDate = pgtype.Date{Time: parsed, Valid: true}
	}

	updated, err := q.UpdateVictim(r.Context(), sqlc.UpdateVictimParams{
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		AddressID: existing.AddressID,
		VictimID:  victimID,
	})
	if err != nil {
		http.Error(w, "failed to update victim", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(victimPartialResponse{
		VictimID:      updated.VictimID,
		MobileUserID:  updated.MobileUserID,
		FirstName:     updated.FirstName,
		LastName:      updated.LastName,
		MedicalInfoID: updated.MedicalInfoID,
	})
}

func GetVictimDetails(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	victimID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("GetVictimDetails [1]: invalid victim id %q: %v", idStr, err)
		http.Error(w, "invalid victim id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		log.Printf("GetVictimDetails [2]: database connection error (victim_id=%d): %v", victimID, err)
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	details, err := q.GetVictimDetails(r.Context(), victimID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "victim not found", http.StatusNotFound)
			return
		}
		log.Printf("GetVictimDetails [3]: query error (victim_id=%d): %v", victimID, err)
		http.Error(w, "failed to fetch victim details", http.StatusInternalServerError)
		return
	}

	birthDate := dateValue(details.BirthDate)
	age := 0
	if details.BirthDate.Valid {
		age = calculateAge(details.BirthDate.Time)
	}

	primaryContactText := formatPrimaryContact(
		textValue(details.PrimaryContactPhone),
		textValue(details.PrimaryContactRelationship),
	)

	response := VictimDetailsResponse{
		VictimID:         details.VictimID,
		FirstName:        details.FirstName,
		LastName:         details.LastName,
		FullName:         fmt.Sprintf("%s %s", details.FirstName, details.LastName),
		BirthDate:        birthDate,
		Age:              age,
		PrimaryContact:   primaryContactText,
		Address:          interfaceString(details.AddressText),
		PregnancyStatus:  textValue(details.PregnancyStatusName),
		OrganDonor:       textValue(details.DonorStatusName),
		BloodType:        textValue(details.BloodTypeName),
		LastDiagnosis:    interfaceString(details.LastDiagnosis),
		DiagnosisDate:    dateValue(details.LastDiagnosisDate),
		PlaceOfDiagnosis: textValue(details.LastDiagnosisHospitalName),
		Allergies:        textValue(details.Allergies),
		Medication:       textValue(details.Medication),
		MedicalHistory:   interfaceString(details.MedicalHistory),
		MedicalNote:      textValue(details.MedicalNotes),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func dateValue(value pgtype.Date) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format("2006-01-02")
}

func bytesValue(value []byte) string {
	if len(value) == 0 {
		return ""
	}
	return string(value)
}

func formatPrimaryContact(phoneNumber, relationship string) string {
	if phoneNumber == "" && relationship == "" {
		return ""
	}
	if phoneNumber == "" {
		return relationship
	}
	if relationship == "" {
		return phoneNumber
	}
	return fmt.Sprintf("%s (%s)", phoneNumber, relationship)
}

func calculateAge(birthDate time.Time) int {
	now := time.Now().UTC()
	years := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}
