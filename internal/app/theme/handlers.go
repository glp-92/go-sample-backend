package theme

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func CreateThemeHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
	var request CreateThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	theme, err := service.CreateTheme(request)
	if err != nil {
		http.Error(w, "Error al crear tema", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateThemeResponse{
		ThemeID:       theme.Id.String(),
		Name:          theme.Name,
		Slug:          theme.Slug,
		Excerpt:       theme.Excerpt,
		FeaturedImage: theme.FeaturedImage,
	})
}

func GetThemeByIDHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	theme, err := service.FindThemeById(id)
	if err != nil {
		http.Error(w, "Error al buscar el tema", http.StatusNotFound)
		return
	}
	if theme.Id == uuid.Nil {
		http.Error(w, "Tema no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(theme)
}

func DeleteThemeHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	err = service.DeleteThemeById(id)
	if err != nil {
		http.Error(w, "Error al eliminar tema", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateThemeByIDHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	var request UpdateThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	theme, err := service.UpdateThemeById(request, id)
	if err != nil {
		http.Error(w, "Error al actualizar tema", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UpdateThemeResponse{
		Name:          theme.Name,
		Slug:          theme.Slug,
		Excerpt:       theme.Excerpt,
		FeaturedImage: theme.FeaturedImage,
	})
}
