package post

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type PostRepository interface {
	Save(post Post) error
	FindByID(id uuid.UUID) (*Post, error)
}

type MySQLPostRepository struct {
	db *sql.DB
}

func NewMySQLPostRepository(db *sql.DB) PostRepository {
	return &MySQLPostRepository{db: db}
}

func (r *MySQLPostRepository) Save(post Post) error {
	query := `
		INSERT INTO posts (id, title, slug, excerpt, content, date)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, post.Id, post.Title, post.Slug, post.Excerpt, post.Content, post.Date)
	return err
}

func (r *MySQLPostRepository) FindByID(id uuid.UUID) (*Post, error) {
	query := `
        SELECT id, title, slug, excerpt, content, date
        FROM posts
        WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var post Post
	err := row.Scan(&post.Id, &post.Title, &post.Slug, &post.Excerpt, &post.Content, &post.Date)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // No encontrado
	}
	if err != nil {
		return nil, err // Error inesperado
	}

	return &post, nil
}
