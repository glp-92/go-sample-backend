package category

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Save(category Category) error
	FindByID(id uuid.UUID) (*Category, error)
}

type MySQLCategoryRepository struct {
	db *sql.DB
}

func NewMySQLCategoryRepository(db *sql.DB) CategoryRepository {
	return &MySQLCategoryRepository{db: db}
}

func (r *MySQLCategoryRepository) Save(category Category) error {
	query := `
		INSERT INTO categories (id, name, slug)
		VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, category.Id, category.Name, category.Slug)
	return err
}

func (r *MySQLCategoryRepository) FindByID(id uuid.UUID) (*Category, error) {
	query := `
        SELECT id, name, slug
        FROM categories
        WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var category Category
	err := row.Scan(&category.Id, &category.Name, &category.Slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}
