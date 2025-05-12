package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthMiddlewares struct {
	authService *AuthService
}

func NewAuthMiddlewares(authService *AuthService) *AuthMiddlewares {
	return &AuthMiddlewares{authService: authService}
}

func GetUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

func (m *AuthMiddlewares) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStrParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenStrParts) != 2 || strings.ToLower(tokenStrParts[0]) != "bearer" {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		access_token := tokenStrParts[1]
		user, err := m.authService.ValidateTokenFromUser(access_token)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		next.ServeHTTP(writer, req.WithContext(ctx))
	})
}

func (m *AuthMiddlewares) Expired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStrParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenStrParts) != 2 || strings.ToLower(tokenStrParts[0]) != "bearer" {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		access_token := tokenStrParts[1]
		user, err := m.authService.ValidateExpiredTokenFromUser(access_token)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		next.ServeHTTP(writer, req.WithContext(ctx))
	})
}
