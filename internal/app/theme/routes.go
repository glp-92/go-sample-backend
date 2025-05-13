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
	http.HandleFunc("GET /themes/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetThemeByIDHandler(service, w, r)
	})
}
