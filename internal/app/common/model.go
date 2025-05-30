package common

import (
	"time"

	"github.com/google/uuid"

	"fullstackcms/backend/internal/app/category"
	"fullstackcms/backend/internal/app/theme"
)

type PostAggregated struct {
	Id            uuid.UUID                          `json:"id"`
	Title         string                             `json:"title"`
	Slug          string                             `json:"slug"`
	Excerpt       string                             `json:"excerpt"`
	FeaturedImage string                             `json:"featuredImage"`
	UserId        uuid.UUID                          `json:"userId"`
	Date          time.Time                          `json:"date"`
	Categories    []category.CategoryDetailsResponse `json:"categories"`
	Themes        []theme.ThemeBasicInfoResponse     `json:"themes"`
}
