package category

import (
	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoriesPageableResponse struct {
	Categories []CategoryDetailsResponse `json:"categories"`
	Total      int                       `json:"totalCategories"`
	Page       int                       `json:"page"`
	PerPage    int                       `json:"perPage"`
	Pages      int                       `json:"pages"`
}

type CategoriesListResponse struct {
	Categories []CategoryDetailsResponse `json:"categories"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CreateCategoryResponse struct {
	CategoryID string `json:"categoryId"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
}

type UpdateCategoryResponse struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoryDetailsResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}
