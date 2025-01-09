package post

import (
	"database/sql"
	"net/http"
)

func RegisterRoutes(db *sql.DB) {
	repo := NewMySQLPostRepository(db)
	service := NewPostService(repo)
	http.HandleFunc("POST /posts", func(w http.ResponseWriter, r *http.Request) {
		CreatePostHandler(service, w, r)
	})
	http.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetPostByIDHandler(service, w, r)
	})
}
