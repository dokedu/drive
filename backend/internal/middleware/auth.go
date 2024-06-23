package middleware

import (
	"example/internal/database"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// AuthMiddleware creates an Auth middleware with the given database connection.
func AuthMiddleware(db *database.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request has the correct token
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := db.FindSessionByToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the session is older than 7 days
			if session.CreatedAt.Before(time.Now().AddDate(0, 0, -7)) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				_, _ = db.RemoveSession(r.Context(), session.ID)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
