package auth

import (
	"database/sql"
	"net/http"
)

func RegisterRoutes(db *sql.DB, service *AuthService, authMiddlewares *AuthMiddlewares) {
	http.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(service, w, r)
	})
	http.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		LoginUserHandler(service, w, r)
	})
	http.Handle("POST /auth/refresh", authMiddlewares.Expired(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RefreshTokenHandler(service, w, r)
	})))
}
