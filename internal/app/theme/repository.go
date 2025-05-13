package theme

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type ThemeRepository interface {
	Save(theme Theme) error
	FindByID(id uuid.UUID) (*Theme, error)
}

type MySQLThemeRepository struct {
	db *sql.DB
}

func NewMySQLThemeRepository(db *sql.DB) ThemeRepository {
	return &MySQLThemeRepository{db: db}
}

func (r *MySQLThemeRepository) Save(theme Theme) error {
	query := `
		INSERT INTO themes (id, name, slug)
		VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, theme.Id, theme.Name, theme.Slug)
	return err
}

func (r *MySQLThemeRepository) FindByID(id uuid.UUID) (*Theme, error) {
	query := `
        SELECT id, name, slug
        FROM themes
        WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var theme Theme
	err := row.Scan(&theme.Id, &theme.Name, &theme.Slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &theme, nil
}
