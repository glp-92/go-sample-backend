package theme

import (
	"github.com/google/uuid"
)

type ThemeService struct {
	repo ThemeRepository
}

func NewThemeService(repo ThemeRepository) *ThemeService {
	return &ThemeService{repo: repo}
}

func (s *ThemeService) CreateTheme(request CreateThemeRequest) (Theme, error) {
	newTheme := Theme{
		Id:            uuid.New(),
		Name:          request.Name,
		Slug:          request.Slug,
		Excerpt:       request.Excerpt,
		FeaturedImage: request.FeaturedImage,
	}
	err := s.repo.Save(newTheme)
	return newTheme, err
}

func (s *ThemeService) FindThemeById(id uuid.UUID) (ThemeDetailsResponse, error) {
	theme, err := s.repo.FindByID(id)
	if err != nil {
		return ThemeDetailsResponse{}, err
	}
	if theme == nil {
		return ThemeDetailsResponse{}, nil
	}
	response := ThemeDetailsResponse{
		Id:            theme.Id,
		Name:          theme.Name,
		Slug:          theme.Slug,
		Excerpt:       theme.Excerpt,
		FeaturedImage: theme.FeaturedImage,
	}
	return response, nil
}

func (s *ThemeService) FindThemesPageable(page, perPage int, reverse bool, queryLen int) ([]Theme, int, error) {
	if queryLen > 1 {
		return s.repo.FindPageable(page, perPage, reverse)
	} else {
		return s.repo.FindAll()
	}
}

func (s *ThemeService) DeleteThemeById(id uuid.UUID) error {
	err := s.repo.DeleteById(id)
	return err
}

func (s *ThemeService) UpdateThemeById(request UpdateThemeRequest, id uuid.UUID) (Theme, error) {
	updatedCategory := Theme{
		Id:            id,
		Name:          request.Name,
		Slug:          request.Slug,
		Excerpt:       request.Excerpt,
		FeaturedImage: request.FeaturedImage,
	}
	err := s.repo.Update(updatedCategory)
	return updatedCategory, err
}
