package post

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type PostRepository interface {
	Save(post Post, categoryIds []uuid.UUID, themeIds []uuid.UUID) error
	FindByID(id uuid.UUID) (*Post, error)
}

type MySQLPostRepository struct {
	db *sql.DB
}

func NewMySQLPostRepository(db *sql.DB) PostRepository {
	return &MySQLPostRepository{db: db}
}

func (r *MySQLPostRepository) Save(post Post, categoryIds []uuid.UUID, themeIds []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO posts (id, title, slug, excerpt, content, featured_image, user_id, date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, post.Id, post.Title, post.Slug, post.Excerpt, post.Content, post.FeaturedImage, post.UserId, post.Date)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(categoryIds) > 0 {
		values := make([]interface{}, 0, len(categoryIds)*2)
		placeholders := make([]string, 0, len(categoryIds))

		for _, categoryId := range categoryIds {
			placeholders = append(placeholders, "(?, ?)")
			values = append(values, post.Id, categoryId)
		}

		catQuery := fmt.Sprintf(`
			INSERT INTO posts_categories (post_id, category_id)
			VALUES %s`, strings.Join(placeholders, ", "))

		_, err = tx.Exec(catQuery, values...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(themeIds) > 0 {
		values := make([]interface{}, 0, len(themeIds)*2)
		placeholders := make([]string, 0, len(themeIds))

		for _, themeId := range themeIds {
			placeholders = append(placeholders, "(?, ?)")
			values = append(values, post.Id, themeId)
		}

		themeQuery := fmt.Sprintf(`
			INSERT INTO posts_themes (post_id, theme_id)
			VALUES %s`, strings.Join(placeholders, ", "))

		_, err = tx.Exec(themeQuery, values...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *MySQLPostRepository) FindByID(id uuid.UUID) (*Post, error) {
	query := `
        SELECT id, title, slug, excerpt, content, featured_image, user_id, date
        FROM posts
        WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var post Post
	err := row.Scan(&post.Id, &post.Title, &post.Slug, &post.Excerpt, &post.Content, &post.FeaturedImage, &post.UserId, &post.Date)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}
