package category

import (
	"fmt"

	"github.com/google/uuid"
)

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(request CreateCategoryRequest, userId uuid.UUID) (Category, error) {
	fmt.Println(userId)
	newCategory := Category{
		Id:   uuid.New(),
		Name: request.Name,
		Slug: request.Slug,
	}
	err := s.repo.Save(newCategory)
	return newCategory, err
}

func (s *CategoryService) FindCategoryById(id uuid.UUID) (CategoryDetailsResponse, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return CategoryDetailsResponse{}, err
	}
	if category == nil {
		return CategoryDetailsResponse{}, nil
	}
	response := CategoryDetailsResponse{
		Id:   category.Id,
		Name: category.Name,
		Slug: category.Slug,
	}
	return response, nil
}
