package category

import (
	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CreateCategoryResponse struct {
	CategoryID string `json:"categoryId"`
}

type CategoryDetailsResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}
