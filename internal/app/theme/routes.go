package theme

import (
	"database/sql"
	"net/http"

	"fullstackcms/backend/pkg/auth"
)

func RegisterRoutes(db *sql.DB, authMiddlewares *auth.AuthMiddlewares) {
	repo := NewMySQLThemeRepository(db)
	service := NewThemeService(repo)
	http.Handle("POST /themes", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateThemeHandler(service, w, r)
	})))
	http.HandleFunc("GET /themes", func(w http.ResponseWriter, r *http.Request) {
		GetThemesWithFiltersHandler(service, w, r)
	})
	http.HandleFunc("GET /themes/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetThemeByIdHandler(service, w, r)
	})
	http.Handle("PUT /themes/{id}", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UpdateThemeByIdHandler(service, w, r)
	})))
	http.Handle("DELETE /themes/{id}", authMiddlewares.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		DeleteThemeByIdHandler(service, w, r)
	})))
}
