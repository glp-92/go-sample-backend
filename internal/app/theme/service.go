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

func (s *ThemeService) CreateTheme(request CreateThemeRequest, userId uuid.UUID) (Theme, error) {
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

func (s *ThemeService) DeleteThemeById(id uuid.UUID) error {
	err := s.repo.DeleteById(id)
	return err
}
