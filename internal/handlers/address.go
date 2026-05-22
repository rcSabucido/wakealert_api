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

type regionResponse struct {
	RegionPsgc string `json:"region_psgc"`
	RegionName string `json:"region_name"`
}

type provinceOrHucResponse struct {
	ProvinceOrHucPsgc string `json:"province_or_huc_psgc"`
	ProvinceOrHucName string `json:"province_or_huc_name"`
	RegionPsgc        string `json:"region_psgc"`
}

type municipalOrCityResponse struct {
	CityMunPsgc       string `json:"city_mun_psgc"`
	CityMunName       string `json:"city_mun_name"`
	ProvinceOrHucPsgc string `json:"province_or_huc_psgc"`
	Type              string `json:"type"`
}

type barangayResponse struct {
	BarangayPsgc string `json:"barangay_psgc"`
	BarangayName string `json:"barangay_name"`
	CityMunPsgc  string `json:"city_mun_psgc"`
}

type addressLineRequest struct {
	BarangayPsgc string `json:"barangay_psgc"`
	AddressLine  string `json:"address_line"`
}

type addressLineResponse struct {
	AddressID    int64   `json:"address_id"`
	AddressLine  string  `json:"address_line"`
	BarangayPsgc *string `json:"barangay_psgc"`
}

type getAddressLineResponse struct {
	AddressID   int64   `json:"address_id"`
	BarangayID  *string `json:"barangay_id"`
	AddressLine string  `json:"address_line"`
}

func ListRegions(w http.ResponseWriter, r *http.Request) {
	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	regions, err := q.ListRegions(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch regions", http.StatusInternalServerError)
		return
	}

	response := make([]regionResponse, 0, len(regions))
	for _, region := range regions {
		response = append(response, regionResponse{
			RegionPsgc: region.RegionPsgc,
			RegionName: region.RegionName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetProvinceOrHucsByRegion(w http.ResponseWriter, r *http.Request) {
	regionPsgc := r.PathValue("region_psgc")
	if regionPsgc == "" {
		http.Error(w, "invalid region psgc", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	provinces, err := q.GetProvinceOrHucsByRegion(r.Context(), regionPsgc)
	if err != nil {
		http.Error(w, "failed to fetch provinces", http.StatusInternalServerError)
		return
	}

	response := make([]provinceOrHucResponse, 0, len(provinces))
	for _, province := range provinces {
		response = append(response, provinceOrHucResponse{
			ProvinceOrHucPsgc: province.ProvinceOrHucPsgc,
			ProvinceOrHucName: province.ProvinceOrHucName,
			RegionPsgc:        province.RegionPsgc,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMunicipalOrCitiesByProvinceOrHuc(w http.ResponseWriter, r *http.Request) {
	provinceOrHucPsgc := r.PathValue("province_or_huc_psgc")
	if provinceOrHucPsgc == "" {
		http.Error(w, "invalid province or huc psgc", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	cities, err := q.GetMunicipalOrCitiesByProvinceOrHuc(r.Context(), pgtype.Text{String: provinceOrHucPsgc, Valid: true})
	if err != nil {
		http.Error(w, "failed to fetch cities", http.StatusInternalServerError)
		return
	}

	response := make([]municipalOrCityResponse, 0, len(cities))
	for _, city := range cities {
		response = append(response, municipalOrCityResponse{
			CityMunPsgc:       city.CityMunPsgc,
			CityMunName:       city.CityMunName,
			ProvinceOrHucPsgc: textValue(city.ProvinceOrHucPsgc),
			Type:              textValue(city.Type),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetBarangaysByMunicipalOrCity(w http.ResponseWriter, r *http.Request) {
	cityMunPsgc := r.PathValue("city_mun_psgc")
	if cityMunPsgc == "" {
		http.Error(w, "invalid city or municipality psgc", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	barangays, err := q.GetBarangaysByMunicipalOrCity(r.Context(), cityMunPsgc)
	if err != nil {
		http.Error(w, "failed to fetch barangays", http.StatusInternalServerError)
		return
	}

	response := make([]barangayResponse, 0, len(barangays))
	for _, barangay := range barangays {
		response = append(response, barangayResponse{
			BarangayPsgc: barangay.BarangayPsgc,
			BarangayName: barangay.BarangayName,
			CityMunPsgc:  barangay.CityMunPsgc,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateAddressLine(w http.ResponseWriter, r *http.Request) {
	var req addressLineRequest

	/*
		if req.BarangayPsgc == "" || req.AddressLine == "" {
			http.Error(w, "barangay_psgc and address_line are required", http.StatusBadRequest)
			return
		}
	*/

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	created, err := q.CreateAddressLine(r.Context(), sqlc.CreateAddressLineParams{
		BarangayPsgc: pgtype.Text{String: req.BarangayPsgc, Valid: req.BarangayPsgc != ""},
		AddressLine:  req.AddressLine,
	})
	if err != nil {
		http.Error(w, "failed to create address line", http.StatusInternalServerError)
		return
	}

	var barangay *string
	if created.BarangayPsgc.Valid {
		barangay = &created.BarangayPsgc.String
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(addressLineResponse{
		AddressID:    created.AddressID,
		AddressLine:  created.AddressLine,
		BarangayPsgc: barangay,
	})
}

func UpdateAddressLine(w http.ResponseWriter, r *http.Request) {
	addressIDStr := r.PathValue("address_id")
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid address id", http.StatusBadRequest)
		return
	}

	var req addressLineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.BarangayPsgc == "" || req.AddressLine == "" {
		http.Error(w, "barangay_psgc and address_line are required", http.StatusBadRequest)
		return
	}

	pool, err := db.Pool(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	var updated addressLineResponse
	if err := pool.QueryRow(r.Context(),
		"UPDATE address_line SET barangay_psgc = $1, address_line = $2 WHERE address_id = $3 RETURNING address_id, address_line, barangay_psgc",
		req.BarangayPsgc,
		req.AddressLine,
		addressID,
	).Scan(&updated.AddressID, &updated.AddressLine, &updated.BarangayPsgc); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "address line not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update address line", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func GetAddressLine(w http.ResponseWriter, r *http.Request) {
	addressIDStr := r.PathValue("address_id")
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid address id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	addressLine, err := q.GetAddressLine(r.Context(), addressID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "address line not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch address line", http.StatusInternalServerError)
		return
	}

	var barangayID *string
	if addressLine.BarangayPsgc.Valid {
		barangayID = &addressLine.BarangayPsgc.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getAddressLineResponse{
		AddressID:   addressLine.AddressID,
		BarangayID:  barangayID,
		AddressLine: addressLine.AddressLine,
	})
}

func textValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
