package theme

import (
	"github.com/google/uuid"
)

type Theme struct {
	Id            uuid.UUID
	Name          string
	Slug          string
	Excerpt       string
	FeaturedImage string
}
