package common

import (
	"database/sql"
	"fmt"
	"fullstackcms/backend/internal/app/category"
	"fullstackcms/backend/internal/app/theme"
	"fullstackcms/backend/pkg/auth"
	"strings"

	"github.com/google/uuid"
)

func FindPostsWithCategoriesAndThemesFiltered(db *sql.DB, keyword, categoryname, themename string, limit, offset int, reverse bool) ([]PostAggregated, int, error) {
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
	if themename != "" {
		conditions = append(conditions, "t.name = ?")
		args = append(args, themename)
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
			GROUP_CONCAT(DISTINCT c.id) as categories_ids,
			GROUP_CONCAT(DISTINCT c.name) as categories_names,
			GROUP_CONCAT(DISTINCT c.slug) as categories_slugs,
			GROUP_CONCAT(DISTINCT t.id) as themes_ids,
			GROUP_CONCAT(DISTINCT t.name) as themes_names,
			GROUP_CONCAT(DISTINCT t.slug) as themes_slugs
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
	var postsAggregated []PostAggregated
	for rows.Next() {
		var p PostAggregated
		var u auth.User
		var categoriesIds sql.NullString
		var categoriesNames sql.NullString
		var categoriesSlugs sql.NullString
		var themesIds sql.NullString
		var themesNames sql.NullString
		var themesSlugs sql.NullString
		err := rows.Scan(
			&p.Id, &p.Title, &p.Slug, &p.Excerpt, &p.FeaturedImage,
			&p.Date, &u.Username, &categoriesIds, &categoriesNames, &categoriesSlugs, &themesIds, &themesNames, &themesSlugs,
		)
		if err != nil {
			return nil, 0, err
		}
		p.Categories = []category.CategoryDetailsResponse{}
		if categoriesIds.Valid && categoriesNames.Valid && categoriesSlugs.Valid &&
			categoriesIds.String != "" && categoriesNames.String != "" && categoriesSlugs.String != "" {
			ids := strings.Split(categoriesIds.String, ",")
			names := strings.Split(categoriesNames.String, ",")
			slugs := strings.Split(categoriesSlugs.String, ",")
			if len(ids) == len(names) && len(ids) == len(slugs) {
				for i := range ids {
					parsedID, err := uuid.Parse(ids[i])
					if err != nil {
						continue
					}
					p.Categories = append(p.Categories, category.CategoryDetailsResponse{
						Id:   parsedID,
						Name: names[i],
						Slug: slugs[i],
					})
				}
			}
		}
		p.Themes = []theme.ThemeBasicInfoResponse{}
		if themesIds.Valid && themesNames.Valid && themesSlugs.Valid &&
			categoriesIds.String != "" && themesNames.String != "" && themesSlugs.String != "" {
			ids := strings.Split(themesIds.String, ",")
			names := strings.Split(themesNames.String, ",")
			slugs := strings.Split(themesNames.String, ",")
			if len(names) == len(slugs) {
				for i := range names {
					parsedID, err := uuid.Parse(ids[i])
					if err != nil {
						continue
					}
					p.Themes = append(p.Themes, theme.ThemeBasicInfoResponse{
						Id:   parsedID,
						Name: names[i],
						Slug: slugs[i],
					})
				}
			}
		}
		postsAggregated = append(postsAggregated, p)
	}
	return postsAggregated, totalPosts, nil
}
