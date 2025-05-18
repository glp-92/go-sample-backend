package theme

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type ThemeRepository interface {
	Save(theme Theme) error
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
		INSERT INTO themes (id, name, slug, excerpt, featuredImage)
		VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, theme.Id, theme.Name, theme.Slug, theme.Excerpt, theme.FeaturedImage)
	return err
}

func (r *MySQLThemeRepository) FindByID(id uuid.UUID) (*Theme, error) {
	query := `
        SELECT id, name, slug, excerpt, featuredImage
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
	_, err := r.db.Exec(query, id)
	return err
}
