package post

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Title         string      `json:"title"`
	Slug          string      `json:"slug"`
	Excerpt       string      `json:"excerpt"`
	Content       string      `json:"content"`
	FeaturedImage string      `json:"featuredImage"`
	CategoryIds   []uuid.UUID `json:"categoryIds"`
	ThemeIds      []uuid.UUID `json:"themeIds"`
}

type UpdatePostRequest struct {
	Title         string      `json:"title"`
	Slug          string      `json:"slug"`
	Excerpt       string      `json:"excerpt"`
	Content       string      `json:"content"`
	FeaturedImage string      `json:"featuredImage"`
	CategoryIds   []uuid.UUID `json:"categoryIds"`
	ThemeIds      []uuid.UUID `json:"themeIds"`
}

type CreatePostResponse struct {
	PostID string `json:"postId"`
}

type UpdatePostResponse struct {
	PostID string `json:"postId"`
}

type PostDetailsResponse struct {
	Id            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Excerpt       string    `json:"excerpt"`
	Content       string    `json:"content"`
	FeaturedImage string    `json:"featuredImage"`
	UserId        uuid.UUID `json:"userId"`
	Date          time.Time `json:"date"`
}

type PostsFilteredResponse struct {
	Posts   []Post `json:"posts"`
	Total   int    `json:"totalPosts"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
	Pages   int    `json:"pages"`
}
