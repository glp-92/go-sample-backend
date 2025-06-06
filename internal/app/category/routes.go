package category

import (
	"database/sql"
	"net/http"

	"fullstackcms/backend/pkg/auth"
)

func RegisterRoutes(db *sql.DB, authMiddlewares *auth.AuthMiddlewares) {
	repo := NewMySQLCategoryRepository(db)
	service := NewCategoryService(repo)
	http.Handle("POST /categories", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateCategoryHandler(service, w, r)
	})))
	http.HandleFunc("GET /categories", func(w http.ResponseWriter, r *http.Request) {
		GetCategoriesWithFiltersHandler(service, w, r)
	})
	http.HandleFunc("GET /categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetCategoryByIdHandler(service, w, r)
	})
	http.Handle("PUT /categories/{id}", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UpdateCategoryByIdHandler(service, w, r)
	})))
	http.Handle("DELETE /categories/{id}", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		DeleteCategoryByIdHandler(service, w, r)
	})))
}
