package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rcsabucido/wakealert_api/internal/db"
	"github.com/rcsabucido/wakealert_api/internal/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type createMobileUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type mobileLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type mobileLoginUserResponse struct {
	MobileUserID int64  `json:"mobile_user_id"`
	Email        string `json:"email"`
}

type mobileLoginResponse struct {
	Token string                  `json:"token"`
	User  mobileLoginUserResponse `json:"user"`
}

type mobileUserResponse struct {
	MobileUserID int64  `json:"mobile_user_id"`
	Email        string `json:"email"`
}

func CreateMobileUser(w http.ResponseWriter, r *http.Request) {
	var req createMobileUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	// Check if a user with the same email already exists
	if existing, err := q.GetMobileUserByEmail(r.Context(), req.Email); err == nil && existing.MobileUserID != 0 {
		http.Error(w, "email already registered", http.StatusConflict)
		return
	} else if err != nil && err != pgx.ErrNoRows {
		http.Error(w, "failed to check existing user", http.StatusInternalServerError)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	mu, err := q.CreateMobileUser(r.Context(), sqlc.CreateMobileUserParams{
		Email:    req.Email,
		Password: string(hashed),
	})
	if err != nil {
		http.Error(w, "failed to create mobile user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mobileUserResponse{
		MobileUserID: mu.MobileUserID,
		Email:        mu.Email,
	})
}

func GetMobileUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid mobile user id", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	mu, err := q.GetMobileUser(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "mobile user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch mobile user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mobileUserResponse{
		MobileUserID: mu.MobileUserID,
		Email:        mu.Email,
	})
}

func GetMobileUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	mu, err := q.GetMobileUserByEmail(r.Context(), email)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "mobile user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch mobile user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mobileUserResponse{
		MobileUserID: mu.MobileUserID,
		Email:        mu.Email,
	})
}

func MobileLogin(w http.ResponseWriter, r *http.Request) {
	var req mobileLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	user, err := q.GetMobileUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.MobileUserID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mobileLoginResponse{
		Token: tokenString,
		User: mobileLoginUserResponse{
			MobileUserID: user.MobileUserID,
			Email:        user.Email,
		},
	})
}
