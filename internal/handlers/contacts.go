package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/rcsabucido/wakealert_api/internal/db"
	"github.com/rcsabucido/wakealert_api/internal/db/sqlc"
)

type contactRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhoneNumber  string `json:"phone_number"`
	ClientUserID int64  `json:"client_user_id"`
	Relationship string `json:"relationship"`
	IsPrimary    bool   `json:"is_primary"`
}

type contactResponse struct {
	ContactID    int64  `json:"contact_id"`
	MobileUserID int64  `json:"mobile_user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhoneNumber  string `json:"phone_number"`
	Relationship string `json:"relationship"`
	IsPrimary    bool   `json:"is_primary"`
	IsDeleted    bool   `json:"is_deleted"`
}

type updateContactRequest struct {
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	PhoneNumber  *string `json:"phone_number"`
	Relationship *string `json:"relationship"`
	IsPrimary    *bool   `json:"is_primary"`
}

type deleteContactByDetailsRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhoneNumber  string `json:"phone_number"`
	ClientUserID int64  `json:"client_user_id"`
	Relationship string `json:"relationship"`
}

type contactDetails struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	PhoneNumber      string `json:"phone_number"`
	RelationshipName string `json:"relationship"`
	IsPrimary        bool   `json:"is_primary"`
}

type updateContactByDetailsRequest struct {
	OriginalContact contactDetails `json:"original_contact"`
	UpdateContact   contactDetails `json:"update_contact"`
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.FirstName == "" || req.LastName == "" || req.PhoneNumber == "" || req.ClientUserID == 0 || req.Relationship == "" {
		http.Error(w, "first_name, last_name, phone_number, client_user_id, and relationship are required", http.StatusBadRequest)
		return
	}

	pool, err := db.Pool(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	var relationshipTypeID int64
	err = pool.QueryRow(r.Context(),
		"SELECT relationship_type_id FROM relationship_type WHERE relationship_name = $1",
		req.Relationship,
	).Scan(&relationshipTypeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "invalid relationship", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to resolve relationship", http.StatusInternalServerError)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	contact, err := q.CreateContact(r.Context(), sqlc.CreateContactParams{
		MobileUserID:       req.ClientUserID,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		PhoneNumber:        req.PhoneNumber,
		RelationshipTypeID: relationshipTypeID,
		IsPrimary:          req.IsPrimary,
	})
	if err != nil {
		http.Error(w, "failed to create contact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contactResponse{
		ContactID:    contact.ContactID,
		MobileUserID: contact.MobileUserID,
		FirstName:    contact.FirstName,
		LastName:     contact.LastName,
		PhoneNumber:  contact.PhoneNumber,
		Relationship: req.Relationship,
		IsPrimary:    contact.IsPrimary,
		IsDeleted:    contact.IsDeleted,
	})
}

func ListContactsByMobileUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("client_user_id")
	clientUserID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid client user id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	contacts, err := q.ListContactsByMobileUser(r.Context(), clientUserID)
	if err != nil {
		http.Error(w, "failed to fetch contacts", http.StatusInternalServerError)
		return
	}

	relationshipCache := make(map[int64]string)
	responses := make([]contactResponse, 0, len(contacts))
	for _, contact := range contacts {
		relationshipName, ok := relationshipCache[contact.RelationshipTypeID]
		if !ok {
			relationship, err := q.GetRelationshipType(r.Context(), contact.RelationshipTypeID)
			if err != nil {
				if err == pgx.ErrNoRows {
					http.Error(w, "relationship type not found", http.StatusInternalServerError)
					return
				}
				http.Error(w, "failed to resolve relationship", http.StatusInternalServerError)
				return
			}
			relationshipName = relationship.RelationshipName
			relationshipCache[contact.RelationshipTypeID] = relationshipName
		}

		responses = append(responses, contactResponse{
			ContactID:    contact.ContactID,
			MobileUserID: contact.MobileUserID,
			FirstName:    contact.FirstName,
			LastName:     contact.LastName,
			PhoneNumber:  contact.PhoneNumber,
			Relationship: relationshipName,
			IsPrimary:    contact.IsPrimary,
			IsDeleted:    contact.IsDeleted,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	contactIDStr := r.PathValue("contact_id")
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid contact id", http.StatusBadRequest)
		return
	}

	clientUserIDStr := r.PathValue("client_user_id")
	clientUserID, err := strconv.ParseInt(clientUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid client user id", http.StatusBadRequest)
		return
	}

	var req updateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.FirstName == nil && req.LastName == nil && req.PhoneNumber == nil && req.Relationship == nil && req.IsPrimary == nil {
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
	if req.PhoneNumber != nil && *req.PhoneNumber == "" {
		http.Error(w, "phone_number cannot be empty", http.StatusBadRequest)
		return
	}
	if req.Relationship != nil && *req.Relationship == "" {
		http.Error(w, "relationship cannot be empty", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	existing, err := q.GetContact(r.Context(), contactID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "contact not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch contact", http.StatusInternalServerError)
		return
	}
	if existing.MobileUserID != clientUserID {
		http.Error(w, "contact not found", http.StatusNotFound)
		return
	}

	relationshipTypeID := existing.RelationshipTypeID
	if req.Relationship != nil {
		pool, err := db.Pool(r.Context())
		if err != nil {
			http.Error(w, "database connection error", http.StatusInternalServerError)
			return
		}
		err = pool.QueryRow(r.Context(),
			"SELECT relationship_type_id FROM relationship_type WHERE relationship_name = $1",
			*req.Relationship,
		).Scan(&relationshipTypeID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "invalid relationship", http.StatusBadRequest)
				return
			}
			http.Error(w, "failed to resolve relationship", http.StatusInternalServerError)
			return
		}
	}

	firstName := existing.FirstName
	if req.FirstName != nil {
		firstName = *req.FirstName
	}
	lastName := existing.LastName
	if req.LastName != nil {
		lastName = *req.LastName
	}
	phoneNumber := existing.PhoneNumber
	if req.PhoneNumber != nil {
		phoneNumber = *req.PhoneNumber
	}
	isPrimary := existing.IsPrimary
	if req.IsPrimary != nil {
		isPrimary = *req.IsPrimary
	}

	updated, err := q.UpdateContact(r.Context(), sqlc.UpdateContactParams{
		FirstName:          firstName,
		LastName:           lastName,
		PhoneNumber:        phoneNumber,
		RelationshipTypeID: relationshipTypeID,
		IsPrimary:          isPrimary,
		ContactID:          contactID,
	})
	if err != nil {
		http.Error(w, "failed to update contact", http.StatusInternalServerError)
		return
	}

	relationshipName := ""
	if req.Relationship != nil {
		relationshipName = *req.Relationship
	} else {
		relationship, err := q.GetRelationshipType(r.Context(), updated.RelationshipTypeID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "relationship type not found", http.StatusInternalServerError)
				return
			}
			http.Error(w, "failed to resolve relationship", http.StatusInternalServerError)
			return
		}
		relationshipName = relationship.RelationshipName
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contactResponse{
		ContactID:    updated.ContactID,
		MobileUserID: updated.MobileUserID,
		FirstName:    updated.FirstName,
		LastName:     updated.LastName,
		PhoneNumber:  updated.PhoneNumber,
		Relationship: relationshipName,
		IsPrimary:    updated.IsPrimary,
		IsDeleted:    updated.IsDeleted,
	})
}

func ClearPrimaryContactsByMobileUser(w http.ResponseWriter, r *http.Request) {
	clientUserIDStr := r.PathValue("client_user_id")
	clientUserID, err := strconv.ParseInt(clientUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid client user id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	if err := q.ClearPrimaryContactsByMobileUser(r.Context(), clientUserID); err != nil {
		http.Error(w, "failed to clear primary contacts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	contactIDStr := r.PathValue("contact_id")
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid contact id", http.StatusBadRequest)
		return
	}

	clientUserIDStr := r.PathValue("client_user_id")
	clientUserID, err := strconv.ParseInt(clientUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid client user id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	if err := q.DeleteContact(r.Context(), sqlc.DeleteContactParams{
		MobileUserID: clientUserID,
		ContactID:    contactID,
	}); err != nil {
		http.Error(w, "failed to delete contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteContactByDetails(w http.ResponseWriter, r *http.Request) {
	var req deleteContactByDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	if req.FirstName == "" || req.LastName == "" || req.PhoneNumber == "" || req.ClientUserID == 0 || req.Relationship == "" {
		http.Error(w, "first_name, last_name, phone_number, client_user_id, and relationship are required", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	if err := q.DeleteContactByDetails(r.Context(), sqlc.DeleteContactByDetailsParams{
		MobileUserID:     req.ClientUserID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		PhoneNumber:      req.PhoneNumber,
		RelationshipName: req.Relationship,
	}); err != nil {
		http.Error(w, "failed to delete contact", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateContactByDetails(w http.ResponseWriter, r *http.Request) {
	clientUserIDStr := r.PathValue("client_user_id")
	clientUserID, err := strconv.ParseInt(clientUserIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid client user id", http.StatusBadRequest)
		return
	}

	var req updateContactByDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.OriginalContact.FirstName == "" || req.OriginalContact.LastName == "" || req.OriginalContact.PhoneNumber == "" || req.OriginalContact.RelationshipName == "" {
		http.Error(w, "original_contact must include first_name, last_name, phone_number, and relationship_name", http.StatusBadRequest)
		return
	}
	if req.UpdateContact.FirstName == "" || req.UpdateContact.LastName == "" || req.UpdateContact.PhoneNumber == "" || req.UpdateContact.RelationshipName == "" {
		http.Error(w, "update_contact must include first_name, last_name, phone_number, and relationship_name", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	updated, err := q.UpdateContactByDetails(r.Context(), sqlc.UpdateContactByDetailsParams{
		MobileUserID:       clientUserID,
		FirstName:          req.OriginalContact.FirstName,
		LastName:           req.OriginalContact.LastName,
		PhoneNumber:        req.OriginalContact.PhoneNumber,
		RelationshipName:   req.OriginalContact.RelationshipName,
		IsPrimary:          req.OriginalContact.IsPrimary,
		FirstName_2:        req.UpdateContact.FirstName,
		LastName_2:         req.UpdateContact.LastName,
		PhoneNumber_2:      req.UpdateContact.PhoneNumber,
		RelationshipName_2: req.UpdateContact.RelationshipName,
		IsPrimary_2:        req.UpdateContact.IsPrimary,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "contact not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update contact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contactResponse{
		ContactID:    updated.ContactID,
		MobileUserID: updated.MobileUserID,
		FirstName:    updated.FirstName,
		LastName:     updated.LastName,
		PhoneNumber:  updated.PhoneNumber,
		Relationship: req.UpdateContact.RelationshipName,
		IsPrimary:    updated.IsPrimary,
		IsDeleted:    updated.IsDeleted,
	})
}
