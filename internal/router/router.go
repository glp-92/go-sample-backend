package router

import (
	"database/sql"
	"fullstackcms/backend/internal/app/post"
	"fullstackcms/backend/pkg/auth"
)

func SetupRouter(db *sql.DB, authService *auth.AuthService) {
	authMiddlewares := auth.NewAuthMiddlewares(authService)
	post.RegisterRoutes(db, authMiddlewares)
	auth.RegisterRoutes(db, authService, authMiddlewares)
}
