package theme

import (
	"github.com/google/uuid"
)

type CreateThemeRequest struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Excerpt       string `json:"excerpt"`
	FeaturedImage string `json:"featuredImage"`
}

type ThemesPageableResponse struct {
	Themes  []ThemeDetailsResponse `json:"themes"`
	Total   int                    `json:"totalThemes"`
	Page    int                    `json:"page"`
	PerPage int                    `json:"perPage"`
	Pages   int                    `json:"pages"`
}

type ThemesListResponse struct {
	Themes []ThemeDetailsResponse `json:"themes"`
}

type UpdateThemeRequest struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Excerpt       string `json:"excerpt"`
	FeaturedImage string `json:"featuredImage"`
}

type CreateThemeResponse struct {
	ThemeID       string `json:"themeId"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Excerpt       string `json:"excerpt"`
	FeaturedImage string `json:"featuredImage"`
}

type UpdateThemeResponse struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Excerpt       string `json:"excerpt"`
	FeaturedImage string `json:"featuredImage"`
}

type ThemeDetailsResponse struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Excerpt       string    `json:"excerpt"`
	FeaturedImage string    `json:"featuredImage"`
}

type ThemeBasicInfoResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}
