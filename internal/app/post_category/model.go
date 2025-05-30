package post_category

import (
	"time"

	"github.com/google/uuid"

	"fullstackcms/backend/internal/app/category"
)

type PostCategory struct {
	Id            uuid.UUID
	Title         string
	Slug          string
	Excerpt       string
	FeaturedImage string
	UserId        uuid.UUID
	Date          time.Time
	Categories    []category.Category
}
