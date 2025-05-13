package category

import (
	"encoding/json"
	"net/http"

	"fullstackcms/backend/pkg/auth"

	"github.com/google/uuid"
)

func CreateCategoryHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var request CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	category, err := service.CreateCategory(request, user.Id)
	if err != nil {
		http.Error(w, "Error al crear categoria", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateCategoryResponse{
		CategoryID: category.Id.String(),
	})
}

func GetCategoryByIDHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	category, err := service.FindCategoryById(id)
	if err != nil {
		http.Error(w, "Error al buscar la categoria", http.StatusNotFound)
		return
	}
	if category.Id == uuid.Nil {
		http.Error(w, "Categoria no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}
