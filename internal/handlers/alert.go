package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rcsabucido/wakealert_api/internal/db"
	"github.com/rcsabucido/wakealert_api/internal/db/sqlc"
)

func GetAlert(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid alert id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	alert, err := q.GetAlert(r.Context(), id)
	if err != nil {
		http.Error(w, "alert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

func ListAlerts(w http.ResponseWriter, r *http.Request) {
	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	alerts, err := q.ListAlerts(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func ListAlertsByVictim(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("victim_id")
	victimID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid victim id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	alerts, err := q.ListAlertsByVictim(r.Context(), victimID)
	if err != nil {
		http.Error(w, "failed to fetch alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func CreateAlert(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
		VictimID  int64  `json:"victim_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	var lat, lng pgtype.Numeric
	if err := lat.Scan(body.Latitude); err != nil {
		http.Error(w, "invalid latitude", http.StatusBadRequest)
		return
	}
	if err := lng.Scan(body.Longitude); err != nil {
		http.Error(w, "invalid longitude", http.StatusBadRequest)
		return
	}

	alert, err := q.CreateAlert(r.Context(), sqlc.CreateAlertParams{
		Latitude:  lat,
		Longitude: lng,
		VictimID:  body.VictimID,
	})
	if err != nil {
		http.Error(w, "failed to create alert", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(alert)
}

func SoftDeleteAlert(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid alert id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	alert, err := q.SoftDeleteAlert(r.Context(), id)
	if err != nil {
		http.Error(w, "alert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}

func UpdateAlertStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	alertID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid alert id", http.StatusBadRequest)
		return
	}

	var body struct {
		IsCompleted bool `json:"isCompleted"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	alert, err := q.UpdateAlertStatus(r.Context(), sqlc.UpdateAlertStatusParams{
		AlertID:     alertID,
		IsCompleted: pgtype.Bool{Bool: body.IsCompleted, Valid: true},
	})
	if err != nil {
		http.Error(w, "alert not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alert)
}
