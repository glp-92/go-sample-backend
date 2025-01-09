package router

import (
	"database/sql"
	"fullstackcms/backend/internal/app/post"
	"fullstackcms/backend/internal/app/user"
	"fullstackcms/backend/pkg/auth"
)

func SetupRouter(db *sql.DB) {
	post.RegisterRoutes(db)
	user.RegisterRoutes(db)
	auth.RegisterRoutes(db)
}
