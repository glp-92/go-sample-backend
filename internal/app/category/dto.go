package category

import (
	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
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
