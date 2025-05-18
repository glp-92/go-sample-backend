package category

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Save(category Category) error
	FindByID(id uuid.UUID) (*Category, error)
	DeleteById(id uuid.UUID) error
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

func (r *MySQLCategoryRepository) DeleteById(id uuid.UUID) error {
	query := `
        DELETE
        FROM categories
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
		return fmt.Errorf("category not found")
	}
	return nil
}
