package post

import (
	"database/sql"
	"net/http"

	"fullstackcms/backend/pkg/auth"
)

func RegisterRoutes(db *sql.DB, authMiddlewares *auth.AuthMiddlewares) {
	repo := NewMySQLPostRepository(db)
	service := NewPostService(repo)
	http.Handle("POST /posts", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreatePostHandler(service, w, r)
	})))
	http.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		GetPostsWithFiltersHandler(service, w, r)
	})
	http.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetPostByIDHandler(service, w, r)
	})
}
