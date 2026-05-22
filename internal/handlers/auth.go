package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rcsabucido/wakealert_api/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	ReceiverUserID int64 `json:"receiver_user_id"`
	Username string `json:"username"`
}

type loginResponse struct {
	Token string `json:"token"`
	User loginUserResponse `json:"user"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return 
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}
	
	q, err := db.Queries(r.Context())
	if err != nil {
		http.Error(w, "database connection error", http.StatusInternalServerError)
		return
	}

	user, err := q.GetReceiverUserByUsername(r.Context(), req.Username)
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
		"sub": user.ReceiverUserID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return 
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{
		Token: tokenString,
		User: loginUserResponse{
			ReceiverUserID: user.ReceiverUserID,
			Username: user.Username,
		},
	})
}