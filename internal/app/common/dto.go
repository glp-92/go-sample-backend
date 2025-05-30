package common

type PostsWithCategoriesFilteredResponse struct {
	Posts   []PostAggregated `json:"posts"`
	Total   int              `json:"totalPosts"`
	Page    int              `json:"page"`
	PerPage int              `json:"perPage"`
	Pages   int              `json:"pages"`
}
