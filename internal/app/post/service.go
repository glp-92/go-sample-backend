package post

import (
	"fullstackcms/backend/internal/app/common"
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

func (s *PostService) UpdatePostById(request UpdatePostRequest, userId uuid.UUID, postId uuid.UUID) (Post, error) {
	updatedPost := Post{
		Id:            postId,
		Title:         request.Title,
		Slug:          request.Slug,
		Excerpt:       request.Excerpt,
		Content:       request.Content,
		FeaturedImage: request.FeaturedImage,
		UserId:        userId,
		Date:          time.Now(),
	}
	err := s.repo.Update(updatedPost, request.CategoryIds, request.ThemeIds)
	return updatedPost, err
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

func (s *PostService) FindPostsWithFilters(keyword, category, theme string, page, perPage int, reverse bool) ([]Post, int, error) {
	offset := (page - 1) * perPage
	return s.repo.FindPostsFiltered(keyword, category, theme, perPage, offset, reverse)
}

func (s *PostService) FindPostsWithCategoriesAndThemesFiltered(keyword, category, theme string, page, perPage int, reverse bool) ([]common.PostSummaryAggregated, int, error) {
	offset := (page - 1) * perPage
	return s.repo.FindPostsWithCategoriesAndThemesFiltered(keyword, category, theme, perPage, offset, reverse)
}

func (s *PostService) FindPostDetailsBySlug(slugStr string) (*common.PostDetailsAggregated, error) {
	return s.repo.FindPostDetailsBySlug(slugStr)
}

func (s *PostService) DeletePostById(id uuid.UUID) error {
	err := s.repo.DeleteById(id)
	return err
}
