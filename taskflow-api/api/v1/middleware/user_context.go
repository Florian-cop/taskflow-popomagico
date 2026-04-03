package middleware

import (
	"context"
	"net/http"

	auditInfra "taskflow-api/internal/audit/infrastructure"
)

// UserContext extrait le header X-User-Id et l'injecte dans le context.
func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-Id")
		if userID == "" {
			userID = "default-user"
		}
		ctx := context.WithValue(r.Context(), auditInfra.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
