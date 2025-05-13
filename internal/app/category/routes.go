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
	http.HandleFunc("GET /categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetCategoryByIDHandler(service, w, r)
	})
}
