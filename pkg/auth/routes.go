package auth

import (
	"database/sql"
	"net/http"
)

func RegisterRoutes(db *sql.DB) {
	repo := NewMySQLAuthRepository(db)
	service := NewAuthService(repo)
	http.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(service, w, r)
	})
	http.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		LoginUserHandler(service, w, r)
	})
	http.HandleFunc("POST /auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		RefreshTokenHandler(service, w, r)
	})
}
