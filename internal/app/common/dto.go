package common

import (
	"fullstackcms/backend/internal/app/category"
	"fullstackcms/backend/internal/app/theme"
	"time"

	"github.com/google/uuid"
)

type PostSummaryListResponse struct {
	Posts   []PostSummaryAggregated `json:"posts"`
	Total   int                     `json:"totalPosts"`
	Page    int                     `json:"page"`
	PerPage int                     `json:"perPage"`
	Pages   int                     `json:"pages"`
}

type PostDetailsResponse struct {
	Id            uuid.UUID                          `json:"id"`
	Title         string                             `json:"title"`
	Slug          string                             `json:"slug"`
	Excerpt       string                             `json:"excerpt"`
	Content       string                             `json:"content"`
	FeaturedImage string                             `json:"featuredImage"`
	UserId        uuid.UUID                          `json:"userId"`
	Date          time.Time                          `json:"date"`
	Categories    []category.CategoryDetailsResponse `json:"categories"`
	Themes        []theme.ThemeBasicInfoResponse     `json:"themes"`
}
