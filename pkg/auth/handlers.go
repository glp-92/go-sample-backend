package auth

import (
	"encoding/json"
	"net/http"
)

func CreateUserHandler(service *UserService, w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
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
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	err := service.ValidateUser(request)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
