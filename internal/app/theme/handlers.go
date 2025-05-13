package theme

import (
	"encoding/json"
	"net/http"

	"fullstackcms/backend/pkg/auth"

	"github.com/google/uuid"
)

func CreateThemeHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var request CreateThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	theme, err := service.CreateTheme(request, user.Id)
	if err != nil {
		http.Error(w, "Error al crear tema", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateThemeResponse{
		ThemeID: theme.Id.String(),
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
