package post

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Title   string
	Slug    string
	Excerpt string
	Content string
}

type PostDetailsResponse struct {
	Id      uuid.UUID
	Title   string
	Slug    string
	Excerpt string
	Content string
	Date    time.Time
}
