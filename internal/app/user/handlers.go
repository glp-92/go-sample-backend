package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func CreateUserHandler(service *UserService, w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	err := service.CreateUser(request)
	if err != nil {
		http.Error(w, "Error al crear el post", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetUserByIDHandler(service *UserService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	user, err := service.FindUserById(id)
	if err != nil {
		http.Error(w, "Error al buscar el usuario", http.StatusNotFound)
		return
	}
	if user.Id == uuid.Nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
