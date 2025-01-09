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

func (s *PostService) CreatePost(request CreatePostRequest) error {
	newPost := Post{
		Id:      uuid.New(),
		Title:   request.Title,
		Slug:    request.Slug,
		Excerpt: request.Excerpt,
		Content: request.Content,
		Date:    time.Now(),
	}
	return s.repo.Save(newPost)
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
		Date:    post.Date,
	}
	return response, nil
}
