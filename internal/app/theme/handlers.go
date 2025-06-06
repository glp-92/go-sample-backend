package theme

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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

func GetThemesWithFiltersHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := 1
	perPage := 20
	reverse := query.Get("reverse") == "true"
	if pageStr := query.Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			http.Error(w, "invalid page parameter", http.StatusBadRequest)
			return
		}
		page = p
	}
	if perPageStr := query.Get("perpage"); perPageStr != "" {
		pp, err := strconv.Atoi(perPageStr)
		if err != nil || pp < 1 {
			http.Error(w, "invalid perpage parameter", http.StatusBadRequest)
			return
		}
		perPage = pp
	}
	themes, totalThemes, err := service.FindThemesPageable(page, perPage, reverse, len(query))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "no themes found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error fetching themes", http.StatusBadRequest)
		return
	}
	totalPages := (totalThemes + perPage - 1) / perPage
	if totalPages == 0 && totalThemes > 0 {
		totalPages = 1
	}
	responseThemes := []ThemeDetailsResponse{}
	for _, th := range themes {
		responseThemes = append(responseThemes, ThemeDetailsResponse(th))
	}
	response := ThemesPageableResponse{
		Themes:  responseThemes,
		Total:   totalThemes,
		Page:    page,
		PerPage: perPage,
		Pages:   totalPages,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetThemeByIdHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
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

func DeleteThemeByIdHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
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

func UpdateThemeByIdHandler(service *ThemeService, w http.ResponseWriter, r *http.Request) {
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
