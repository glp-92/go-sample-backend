package post

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"fullstackcms/backend/internal/app/common"
	"fullstackcms/backend/pkg/auth"

	"github.com/google/uuid"
)

func CreatePostHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var request CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	post, err := service.CreatePost(request, user.Id)
	if err != nil {
		http.Error(w, "Error al crear el post", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatePostResponse{
		PostID: post.Id.String(),
	})
}

func UpdatePostByIdHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	postIdStr := r.PathValue("id")
	postId, err := uuid.Parse(postIdStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	var request UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error en el formato del cuerpo", http.StatusBadRequest)
		return
	}
	post, err := service.UpdatePostById(request, user.Id, postId)
	if err != nil {
		http.Error(w, "Error al actualizar el post", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UpdatePostResponse{
		PostID: post.Id.String(),
	})
}

func GetPostByIdHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	post, err := service.FindPostById(id)
	if err != nil {
		http.Error(w, "Error al buscar el post", http.StatusNotFound)
		return
	}
	if post.Id == uuid.Nil {
		http.Error(w, "Post no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func GetPostBySlugHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	slugStr := r.PathValue("slug")
	post, err := service.FindPostDetailsBySlug(slugStr)
	if err != nil || post == nil {
		http.Error(w, "Error al buscar el post", http.StatusNotFound)
		return
	}
	if post.Id == uuid.Nil {
		http.Error(w, "Post no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func GetPostsWithFiltersHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	keyword := query.Get("keyword")
	category := query.Get("category")
	theme := query.Get("theme")
	reverse := query.Get("reverse") == "true"
	pageStr := query.Get("page")
	if pageStr == "" {
		http.Error(w, "missing 'page' parameter", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, "invalid 'page' parameter", http.StatusBadRequest)
		return
	}
	perPageStr := query.Get("perpage")
	perPage := 10
	if perPageStr != "" {
		if p, err := strconv.Atoi(perPageStr); err == nil && p > 0 {
			perPage = p
		}
	}
	posts, totalPosts, err := service.FindPostsWithCategoriesAndThemesFiltered(keyword, category, theme, page, perPage, reverse)
	if err != nil {
		http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if posts == nil {
		posts = []common.PostSummaryAggregated{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(common.PostSummaryListResponse{
		Posts:   posts,
		Total:   totalPosts,
		Page:    page,
		PerPage: perPage,
		Pages:   int(math.Ceil(float64(totalPosts) / float64(perPage))),
	})
}

func DeletePostByIdHandler(service *PostService, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	err = service.DeletePostById(id)
	if err != nil {
		http.Error(w, "Error al eliminar post", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
