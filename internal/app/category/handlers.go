package category

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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

func GetCategoriesWithFiltersHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
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
	categories, totalCategories, err := service.FindCategoriesPageable(page, perPage, reverse, len(query))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "no categories found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error fetching categories", http.StatusBadRequest)
		return
	}
	totalPages := (totalCategories + perPage - 1) / perPage
	if totalPages == 0 && totalCategories > 0 {
		totalPages = 1
	}
	responseCategories := []CategoryDetailsResponse{}
	for _, cat := range categories {
		responseCategories = append(responseCategories, CategoryDetailsResponse(cat))
	}
	response := CategoriesPageableResponse{
		Categories: responseCategories,
		Total:      totalCategories,
		Page:       page,
		PerPage:    perPage,
		Pages:      totalPages,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetCategoryByIdHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
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

func DeleteCategoryByIdHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
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

func UpdateCategoryByIdHandler(service *CategoryService, w http.ResponseWriter, r *http.Request) {
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
