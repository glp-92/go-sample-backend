package router

import (
	"database/sql"
	"fullstackcms/backend/internal/app/category"
	"fullstackcms/backend/internal/app/post"
	"fullstackcms/backend/internal/app/theme"
	"fullstackcms/backend/pkg/auth"
)

func SetupRouter(db *sql.DB, authService *auth.AuthService) {
	authMiddlewares := auth.NewAuthMiddlewares(authService)
	post.RegisterRoutes(db, authMiddlewares)
	category.RegisterRoutes(db, authMiddlewares)
	theme.RegisterRoutes(db, authMiddlewares)
	auth.RegisterRoutes(db, authService, authMiddlewares)
}
