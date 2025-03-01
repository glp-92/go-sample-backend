package auth

import (
	"encoding/json"
	"fmt"
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
	accesToken, refreshToken, err := service.CreateTokens(request, userAgent, user)
	if err != nil {
		http.Error(w, "Token Err", http.StatusUnauthorized)
	}
	refreshTokenCookie := &http.Cookie{Name: "refresh_token", Value: refreshToken, HttpOnly: true}
	http.SetCookie(w, refreshTokenCookie)
	w.Header().Set("Content-Type", "application/json")
	response := dto.LoginResponse{
		AccessToken: accesToken,
	}
	json.NewEncoder(w).Encode(response)
}

func RefreshTokenHandler(service *AuthService, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh Token Err", http.StatusUnauthorized)
	}
	err = service.RefreshToken(cookie.Value)
	if err != nil {
		http.Error(w, "Refresh Token Err", http.StatusUnauthorized)
	}
	fmt.Println(cookie.Value)
}
