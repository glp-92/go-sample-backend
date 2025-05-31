package post

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"fullstackcms/backend/internal/app/common"
	"fullstackcms/backend/pkg/auth"
)

type PostRepository interface {
	Save(post Post, categoryIds []uuid.UUID, themeIds []uuid.UUID) error
	Update(post Post, categoryIds []uuid.UUID, themeIds []uuid.UUID) error
	FindByID(id uuid.UUID) (*Post, error)
	FindPostsFiltered(keyword, category, theme string, limit, offset int, reverse bool) ([]Post, int, error)
	FindPostsWithCategoriesAndThemesFiltered(keyword, category, theme string, limit, offset int, reverse bool) ([]common.PostSummaryAggregated, int, error)
	FindPostDetailsBySlug(slugStr string) (*common.PostDetailsAggregated, error)
	DeleteById(id uuid.UUID) error
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

func (r *MySQLPostRepository) Update(post Post, categoryIds []uuid.UUID, themeIds []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := `
		UPDATE posts
		SET title = ?, slug = ?, excerpt = ?, content = ?, featured_image = ?, user_id = ?, date = ?
		WHERE id = ?`
	_, err = tx.Exec(query, post.Title, post.Slug, post.Excerpt, post.Content, post.FeaturedImage, post.UserId, post.Date, post.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM posts_categories WHERE post_id = ?`, post.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM posts_themes WHERE post_id = ?`, post.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(categoryIds) > 0 {
		values := make([]any, 0, len(categoryIds)*2)
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
		values := make([]any, 0, len(themeIds)*2)
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

func (r *MySQLPostRepository) DeleteById(id uuid.UUID) error {
	query := `
        DELETE
        FROM posts
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

func (r *MySQLPostRepository) FindPostsFiltered(keyword, category, theme string, limit, offset int, reverse bool) ([]Post, int, error) {
	var (
		conditions []string
		args       []any
	)
	if keyword != "" {
		conditions = append(conditions, "(title LIKE ? OR content LIKE ?)")
		kw := "%" + keyword + "%"
		args = append(args, kw, kw)
	}
	if category != "" {
		conditions = append(conditions, "c.name = ?")
		args = append(args, category)
	}
	if theme != "" {
		conditions = append(conditions, "t.name = ?")
		args = append(args, theme)
	}
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.id)
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN posts_categories pc ON p.id = pc.post_id
		LEFT JOIN categories c ON pc.category_id = c.id
		LEFT JOIN posts_themes pt ON p.id = pt.post_id
		LEFT JOIN themes t ON pt.theme_id = t.id
		%s`, whereClause)

	countArgs := make([]any, len(args))
	copy(countArgs, args)
	var totalPosts int
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&totalPosts)
	if err != nil {
		return nil, 0, err
	}
	order := "DESC"
	if reverse {
		order = "ASC"
	}
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
        SELECT p.id, p.title, p.slug, p.excerpt, p.content, p.featured_image, p.date, u.username, GROUP_CONCAT(DISTINCT c.name) as categories_names, GROUP_CONCAT(DISTINCT c.slug) as categories_slugs
        FROM posts p
		JOIN users u on p.user_id = u.id
		LEFT JOIN posts_categories pc ON p.id = pc.post_id
		LEFT JOIN categories c ON pc.category_id = c.id
		LEFT JOIN posts_themes pt ON p.id = pt.post_id
		LEFT JOIN themes t ON pt.theme_id = t.id
        %s
		GROUP BY p.id
        ORDER BY date %s
        LIMIT ? OFFSET ?`, whereClause, order)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		var u auth.User
		var categoriesNames sql.NullString
		var categoriesSlugs sql.NullString
		err := rows.Scan(&p.Id, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.FeaturedImage, &p.Date, &u.Username, &categoriesNames, &categoriesSlugs)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, p)
	}
	return posts, totalPosts, nil
}

func (r *MySQLPostRepository) FindPostsWithCategoriesAndThemesFiltered(keyword, category, theme string, limit, offset int, reverse bool) ([]common.PostSummaryAggregated, int, error) {
	return common.FindPostsWithCategoriesAndThemesFiltered(r.db, keyword, category, theme, limit, offset, reverse)
}

func (r *MySQLPostRepository) FindPostDetailsBySlug(slugStr string) (*common.PostDetailsAggregated, error) {
	return common.FindPostDetailsBySlug(r.db, slugStr)
}
