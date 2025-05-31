package common

import (
	"database/sql"
	"fmt"
	"fullstackcms/backend/internal/app/category"
	"fullstackcms/backend/internal/app/theme"
	"strings"

	"github.com/google/uuid"
)

func FindPostsWithCategoriesAndThemesFiltered(db *sql.DB, keyword, categoryname, themename string, limit, offset int, reverse bool) ([]PostSummaryAggregated, int, error) {
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
		return []PostSummaryAggregated{}, 0, err
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
		return []PostSummaryAggregated{}, 0, err
	}
	defer rows.Close()
	var postsAggregated []PostSummaryAggregated
	for rows.Next() {
		var p PostSummaryAggregated
		var (
			catIdsStr, catNamesStr, catSlugsStr       sql.NullString
			themeIdsStr, themeNamesStr, themeSlugsStr sql.NullString
		)
		err := rows.Scan(
			&p.Id, &p.Title, &p.Slug, &p.Excerpt, &p.FeaturedImage,
			&p.Date, &p.Username, &catIdsStr, &catNamesStr, &catSlugsStr, &themeIdsStr, &themeNamesStr, &themeSlugsStr,
		)
		if err != nil {
			return []PostSummaryAggregated{}, 0, err
		}
		p.Categories = parseCategories(catIdsStr, catNamesStr, catSlugsStr)
		p.Themes = parseThemes(themeIdsStr, themeNamesStr, themeSlugsStr)
		postsAggregated = append(postsAggregated, p)
	}
	return postsAggregated, totalPosts, nil
}

func FindPostDetailsBySlug(db *sql.DB, slugStr string) (*PostDetailsAggregated, error) {
	query := `
        SELECT p.id, p.title, p.slug, p.excerpt, p.content, p.featured_image, p.date, u.username,
            GROUP_CONCAT(DISTINCT c.id),
            GROUP_CONCAT(DISTINCT c.name),
            GROUP_CONCAT(DISTINCT c.slug),
            GROUP_CONCAT(DISTINCT t.id),
            GROUP_CONCAT(DISTINCT t.name),
            GROUP_CONCAT(DISTINCT t.slug)
        FROM posts p
        JOIN users u on p.user_id = u.id
        LEFT JOIN posts_categories pc ON p.id = pc.post_id
        LEFT JOIN categories c ON pc.category_id = c.id
        LEFT JOIN posts_themes pt ON p.id = pt.post_id
        LEFT JOIN themes t ON pt.theme_id = t.id
        WHERE p.slug = ?
        GROUP BY p.id`
	rows, err := db.Query(query, slugStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var postDetails PostDetailsAggregated
	var (
		catIdsStr, catNamesStr, catSlugsStr       sql.NullString
		themeIdsStr, themeNamesStr, themeSlugsStr sql.NullString
	)
	if rows.Next() {
		err := rows.Scan(
			&postDetails.Id,
			&postDetails.Title,
			&postDetails.Slug,
			&postDetails.Excerpt,
			&postDetails.Content,
			&postDetails.FeaturedImage,
			&postDetails.Date,
			&postDetails.Username,
			&catIdsStr,
			&catNamesStr,
			&catSlugsStr,
			&themeIdsStr,
			&themeNamesStr,
			&themeSlugsStr,
		)
		if err != nil {
			return nil, err
		}
		postDetails.Categories = parseCategories(catIdsStr, catNamesStr, catSlugsStr)
		postDetails.Themes = parseThemes(themeIdsStr, themeNamesStr, themeSlugsStr)
		return &postDetails, nil
	}
	return nil, nil
}

func parseCategories(ids, names, slugs sql.NullString) []category.CategoryDetailsResponse {
	if !ids.Valid || !names.Valid || !slugs.Valid {
		return []category.CategoryDetailsResponse{}
	}
	idList := strings.Split(ids.String, ",")
	nameList := strings.Split(names.String, ",")
	slugList := strings.Split(slugs.String, ",")
	var categories []category.CategoryDetailsResponse
	for i := range idList {
		catID, err := uuid.Parse(strings.TrimSpace(idList[i]))
		if err != nil {
			continue
		}
		categories = append(categories, category.CategoryDetailsResponse{
			Id:   catID,
			Name: nameList[i],
			Slug: slugList[i],
		})
	}
	return categories
}

func parseThemes(ids, names, slugs sql.NullString) []theme.ThemeBasicInfoResponse {
	if !ids.Valid || !names.Valid || !slugs.Valid {
		return []theme.ThemeBasicInfoResponse{}
	}
	idList := strings.Split(ids.String, ",")
	nameList := strings.Split(names.String, ",")
	slugList := strings.Split(slugs.String, ",")
	var themes []theme.ThemeBasicInfoResponse
	for i := range idList {
		themeID, err := uuid.Parse(strings.TrimSpace(idList[i]))
		if err != nil {
			continue // salteamos si el UUID no es válido
		}
		themes = append(themes, theme.ThemeBasicInfoResponse{
			Id:   themeID,
			Name: nameList[i],
			Slug: slugList[i],
		})
	}
	return themes
}
