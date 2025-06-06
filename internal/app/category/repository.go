package category

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Save(category Category) error
	Update(Category Category) error
	FindAll() ([]Category, int, error)
	FindPageable(page, perPage int, reverse bool) ([]Category, int, error)
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

func (r *MySQLCategoryRepository) Update(category Category) error {
	query := `
		UPDATE categories
			SET name = ?, slug = ?
			WHERE id = ?`
	_, err := r.db.Exec(query, category.Name, category.Slug, category.Id)
	return err
}

func (r *MySQLCategoryRepository) FindAll() ([]Category, int, error) {
	query := `
        SELECT id, name, slug
        FROM categories
        ORDER BY name ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Name, &category.Slug)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}
	return categories, 0, nil
}

func (r *MySQLCategoryRepository) FindPageable(page, perPage int, reverse bool) ([]Category, int, error) {
	offset := max((page-1)*perPage, 0)
	orderDirection := "ASC"
	if reverse {
		orderDirection = "DESC"
	}
	query := fmt.Sprintf(`
        SELECT id, name, slug
        FROM categories
        ORDER BY name %s
        LIMIT ? OFFSET ?`, orderDirection)
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Name, &category.Slug)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}
	countQuery := `SELECT COUNT(id) FROM categories`
	var totalCategories int
	err = r.db.QueryRow(countQuery).Scan(&totalCategories)
	if err != nil {
		return nil, 0, err
	}
	return categories, totalCategories, nil
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
