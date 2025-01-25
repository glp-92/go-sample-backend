package auth

import (
	"encoding/json"
	"fullstackcms/backend/pkg/auth/dto"
	"net/http"
)

func CreateUserHandler(service *AuthService, w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "DecodeError", http.StatusBadRequest)
		return
	}
	err := service.CreateUser(request)
	if err != nil {
		http.Error(w, "CreateUserError", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginUserHandler(service *AuthService, w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	userAgent := r.Header.Get("User-Agent")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}
	user, err := service.ValidateUser(request)
	if err != nil {
		http.Error(w, "Invalid User", http.StatusUnauthorized)
		return
	}
	tokens, err := service.CreateToken(request, userAgent, user)
	if err != nil {
		http.Error(w, "Token Err", http.StatusUnauthorized)
	}
	cookie1 := &http.Cookie{Name: "refresh_token", Value: tokens.RefreshToken, HttpOnly: true}
	http.SetCookie(w, cookie1)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
