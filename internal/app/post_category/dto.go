package post_category

type PostsWithCategoriesFilteredResponse struct {
	Posts   []PostCategory `json:"posts"`
	Total   int            `json:"totalPosts"`
	Page    int            `json:"page"`
	PerPage int            `json:"perPage"`
	Pages   int            `json:"pages"`
}
