package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rcsabucido/wakealert_api/internal/handlers"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", handlers.Login)
	mux.HandleFunc("POST /mobile_users", handlers.CreateMobileUser)
	mux.HandleFunc("GET /mobile_users/{id}", handlers.GetMobileUser)
	mux.HandleFunc("GET /mobile_users/email/{email}", handlers.GetMobileUserByEmail)
	mux.HandleFunc("POST /medical_info/add", handlers.CreateMedicalInformation)
	mux.HandleFunc("PUT /medical_info/update/{medical_info_id}", handlers.UpdateMedicalInformation)
	mux.HandleFunc("GET /medical_history/{medical_info_id}", handlers.ListMedicalHistoryByMedicalInfo)
	mux.HandleFunc("POST /medical_history", handlers.UpsertMedicalHistoryByMedicalInfo)
	mux.HandleFunc("DELETE /medical_history/{medical_info_id}", handlers.SoftDeleteMedicalHistoryByMedicalInfo)
	mux.HandleFunc("POST /contacts/add", handlers.CreateContact)
	mux.HandleFunc("POST /contacts/edit/{client_user_id}/{contact_id}", handlers.UpdateContact)
	mux.HandleFunc("POST /contacts/edit/{client_user_id}", handlers.UpdateContactByDetails)
	mux.HandleFunc("PUT /contacts/{client_user_id}/primary/clear", handlers.ClearPrimaryContactsByMobileUser)
	mux.HandleFunc("GET /contacts/client/{client_user_id}", handlers.ListContactsByMobileUser)
	mux.HandleFunc("DELETE /contacts/{client_user_id}/{contact_id}", handlers.DeleteContact)
	mux.HandleFunc("DELETE /contacts/by_details", handlers.DeleteContactByDetails)
	mux.HandleFunc("POST /mobile/auth/login", handlers.MobileLogin)
	mux.HandleFunc("GET /alerts/{id}", handlers.GetAlert)
	mux.HandleFunc("PUT /alerts/{id}", handlers.UpdateAlertStatus)
	mux.HandleFunc("DELETE /alerts/{id}", handlers.SoftDeleteAlert)
	mux.HandleFunc("GET /alerts/victim/{victim_id}", handlers.ListAlertsByVictim)
	mux.HandleFunc("GET /alerts", handlers.ListAlerts)
	mux.HandleFunc("POST /alerts", handlers.CreateAlert)
	mux.HandleFunc("GET /victims", handlers.ListVictims)
	mux.HandleFunc("GET /victims/mobile_user/{mobile_user_id}", handlers.GetVictimByMobileUser)
	mux.HandleFunc("GET /victims/{id}", handlers.GetVictimDetails)
	mux.HandleFunc("GET /victims/address_id/{victim_id}", handlers.GetVictimAddressID)
	mux.HandleFunc("PUT /victims/address_id/{victim_id}", handlers.UpdateVictimAddressID)
	mux.HandleFunc("POST /victims/add", handlers.CreateVictim)
	mux.HandleFunc("POST /victims/update/{victim_id}", handlers.UpdateVictim)
	mux.HandleFunc("GET /addresses/regions", handlers.ListRegions)
	mux.HandleFunc("GET /addresses/regions/{region_psgc}/provinces", handlers.GetProvinceOrHucsByRegion)
	mux.HandleFunc("GET /addresses/provinces/{province_or_huc_psgc}/cities", handlers.GetMunicipalOrCitiesByProvinceOrHuc)
	mux.HandleFunc("GET /addresses/cities/{city_mun_psgc}/barangays", handlers.GetBarangaysByMunicipalOrCity)
	mux.HandleFunc("POST /addresses/lines", handlers.CreateAddressLine)
	mux.HandleFunc("PUT /addresses/lines/{address_id}", handlers.UpdateAddressLine)
	mux.HandleFunc("GET /addresses/lines/{address_id}", handlers.GetAddressLine)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}
