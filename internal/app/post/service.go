package post

import (
	"time"

	"github.com/google/uuid"
)

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(request CreatePostRequest, userId uuid.UUID) (Post, error) {
	newPost := Post{
		Id:            uuid.New(),
		Title:         request.Title,
		Slug:          request.Slug,
		Excerpt:       request.Excerpt,
		Content:       request.Content,
		FeaturedImage: request.FeaturedImage,
		UserId:        userId,
		Date:          time.Now(),
	}
	err := s.repo.Save(newPost, request.CategoryIds, request.ThemeIds)
	return newPost, err
}

func (s *PostService) FindPostById(id uuid.UUID) (PostDetailsResponse, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return PostDetailsResponse{}, err
	}
	if post == nil {
		return PostDetailsResponse{}, nil
	}
	response := PostDetailsResponse{
		Id:      post.Id,
		Title:   post.Title,
		Slug:    post.Slug,
		Excerpt: post.Excerpt,
		Content: post.Content,
		UserId:  post.UserId,
		Date:    post.Date,
	}
	return response, nil
}
