package common

import (
	"time"

	"github.com/google/uuid"

	"fullstackcms/backend/internal/app/category"
	"fullstackcms/backend/internal/app/theme"
)

type PostSummaryAggregated struct {
	Id            uuid.UUID                          `json:"id"`
	Title         string                             `json:"title"`
	Slug          string                             `json:"slug"`
	Excerpt       string                             `json:"excerpt"`
	FeaturedImage string                             `json:"featuredImage"`
	Username      string                             `json:"username"`
	Date          time.Time                          `json:"date"`
	Categories    []category.CategoryDetailsResponse `json:"categories"`
	Themes        []theme.ThemeBasicInfoResponse     `json:"themes"`
}

type PostDetailsAggregated struct {
	Id            uuid.UUID                          `json:"id"`
	Title         string                             `json:"title"`
	Slug          string                             `json:"slug"`
	Excerpt       string                             `json:"excerpt"`
	Content       string                             `json:"content"`
	FeaturedImage string                             `json:"featuredImage"`
	Username      string                             `json:"username"`
	Date          time.Time                          `json:"date"`
	Categories    []category.CategoryDetailsResponse `json:"categories"`
	Themes        []theme.ThemeBasicInfoResponse     `json:"themes"`
}
