package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id      uuid.UUID
	Title   string
	Slug    string
	Excerpt string
	Content string
	Date    time.Time
}
