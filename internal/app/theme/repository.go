package theme

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ThemeRepository interface {
	Save(theme Theme) error
	Update(theme Theme) error
	FindByID(id uuid.UUID) (*Theme, error)
	DeleteById(id uuid.UUID) error
}

type MySQLThemeRepository struct {
	db *sql.DB
}

func NewMySQLThemeRepository(db *sql.DB) ThemeRepository {
	return &MySQLThemeRepository{db: db}
}

func (r *MySQLThemeRepository) Save(theme Theme) error {
	query := `
		INSERT INTO themes (id, name, slug, excerpt, featured_image)
		VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, theme.Id, theme.Name, theme.Slug, theme.Excerpt, theme.FeaturedImage)
	return err
}

func (r *MySQLThemeRepository) Update(theme Theme) error {
	query := `
		UPDATE themes
			SET name = ?, slug = ?, excerpt = ?, featured_image = ?
			WHERE id = ?`
	_, err := r.db.Exec(query, theme.Name, theme.Slug, theme.Excerpt, theme.FeaturedImage, theme.Id)
	return err
}

func (r *MySQLThemeRepository) FindByID(id uuid.UUID) (*Theme, error) {
	query := `
        SELECT id, name, slug, excerpt, featured_image
        FROM themes
        WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var theme Theme
	err := row.Scan(&theme.Id, &theme.Name, &theme.Slug, &theme.Excerpt, &theme.FeaturedImage)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *MySQLThemeRepository) DeleteById(id uuid.UUID) error {
	query := `
        DELETE
        FROM themes
        WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("theme not found")
	}
	return nil
}
