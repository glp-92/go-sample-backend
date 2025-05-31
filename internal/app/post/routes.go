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
	http.Handle("PUT /posts/{id}", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UpdatePostByIdHandler(service, w, r)
	})))
	http.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		GetPostsWithFiltersHandler(service, w, r)
	})
	http.HandleFunc("GET /posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		GetPostBySlugHandler(service, w, r)
	})
	http.Handle("DELETE /posts/{id}", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		DeletePostByIdHandler(service, w, r)
	})))
}
