package category

import (
	"github.com/google/uuid"
)

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(request CreateCategoryRequest) (Category, error) {
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

func (s *CategoryService) FindCategoriesPageable(page, perPage int, reverse bool, queryLen int) ([]Category, int, error) {
	if queryLen > 1 {
		return s.repo.FindPageable(page, perPage, reverse)
	} else {
		return s.repo.FindAll()
	}
}

func (s *CategoryService) DeleteCategoryById(id uuid.UUID) error {
	err := s.repo.DeleteById(id)
	return err
}

func (s *CategoryService) UpdateCategoryById(request UpdateCategoryRequest, id uuid.UUID) (Category, error) {
	updatedCategory := Category{
		Id:   id,
		Name: request.Name,
		Slug: request.Slug,
	}
	err := s.repo.Update(updatedCategory)
	return updatedCategory, err
}
