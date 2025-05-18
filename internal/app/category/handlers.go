package category

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func CreateCategoryHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
	var request CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	category, err := service.CreateCategory(request)
	if err != nil {
		http.Error(w, "Error al crear categoria", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateCategoryResponse{
		CategoryID: category.Id.String(),
		Name:       category.Name,
		Slug:       category.Slug,
	})
}

func GetCategoryByIDHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
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

func DeleteCategoryHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	err = service.DeleteCategoryById(id)
	if err != nil {
		http.Error(w, "Error al eliminar categoria", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateCategoryByIDHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	var request UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	category, err := service.UpdateCategoryById(request, id)
	if err != nil {
		http.Error(w, "Error al actualizar categoria", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UpdateCategoryResponse{
		Name: category.Name,
		Slug: category.Slug,
	})
}
