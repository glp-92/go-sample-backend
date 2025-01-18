package auth

import (
	"encoding/json"
	"fullstackcms/backend/pkg/auth/dto"
	"net/http"
)

func CreateUserHandler(service *UserService, w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	err := service.CreateUser(request)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginUserHandler(service *UserService, w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	err := service.ValidateUser(request)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	tokens, err := service.CreateToken(request)
	if err != nil {
		http.Error(w, "Error en login", http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
