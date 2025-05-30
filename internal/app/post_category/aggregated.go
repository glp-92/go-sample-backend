package post_category

import (
	"database/sql"
	"fmt"
	"fullstackcms/backend/internal/app/category"
	"fullstackcms/backend/pkg/auth"
	"strings"
)

func FindPostsWithCategoriesFiltered(db *sql.DB, keyword, categoryname, theme string, limit, offset int, reverse bool) ([]PostCategory, int, error) {
	var (
		conditions []string
		args       []any
	)
	if keyword != "" {
		conditions = append(conditions, "(p.title LIKE ? OR p.content LIKE ?)")
		kw := "%" + keyword + "%"
		args = append(args, kw, kw)
	}
	if categoryname != "" {
		conditions = append(conditions, "c.name = ?")
		args = append(args, categoryname)
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
	err := db.QueryRow(countQuery, countArgs...).Scan(&totalPosts)
	if err != nil {
		return nil, 0, err
	}
	order := "DESC"
	if reverse {
		order = "ASC"
	}
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
        SELECT p.id, p.title, p.slug, p.excerpt, p.featured_image, p.date, u.username,
		       GROUP_CONCAT(DISTINCT c.name) as categories_names,
		       GROUP_CONCAT(DISTINCT c.slug) as categories_slugs
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

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var postsWithCategories []PostCategory
	for rows.Next() {
		var p PostCategory
		var u auth.User
		var categoriesNames sql.NullString
		var categoriesSlugs sql.NullString
		err := rows.Scan(
			&p.Id, &p.Title, &p.Slug, &p.Excerpt, &p.FeaturedImage,
			&p.Date, &u.Username, &categoriesNames, &categoriesSlugs,
		)
		if err != nil {
			return nil, 0, err
		}
		p.Categories = []category.Category{}
		if categoriesNames.Valid && categoriesSlugs.Valid &&
			categoriesNames.String != "" && categoriesSlugs.String != "" {
			names := strings.Split(categoriesNames.String, ",")
			slugs := strings.Split(categoriesSlugs.String, ",")

			if len(names) == len(slugs) {
				for i := range names {
					p.Categories = append(p.Categories, category.Category{
						Name: names[i],
						Slug: slugs[i],
					})
				}
			}
		}
		postsWithCategories = append(postsWithCategories, p)
	}
	return postsWithCategories, totalPosts, nil
}
