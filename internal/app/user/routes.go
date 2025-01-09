package user

import (
	"database/sql"
	"net/http"
)

func RegisterRoutes(db *sql.DB) {
	repo := NewMySQLUserRepository(db)
	service := NewUserService(repo)
	http.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(service, w, r)
	})
	http.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetUserByIDHandler(service, w, r)
	})
}
