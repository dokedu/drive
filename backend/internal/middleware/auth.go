package middleware

import (
	"context"
	"example/internal/database"
	"example/internal/database/db"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ctxKey string

const userCtxKey ctxKey = "user"

// AuthMiddleware creates an Auth middleware with the given database connection.
func AuthMiddleware(dbs *database.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request has the correct token
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := dbs.FindSessionByToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the session is older than 7 days
			if session.CreatedAt.Before(time.Now().AddDate(0, 0, -7)) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				_, _ = dbs.RemoveSession(r.Context(), db.RemoveSessionParams{
					Token:  token,
					UserID: session.UserID,
				})
				return
			}

			user, err := dbs.UserFindByID(r.Context(), session.UserID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add the user to the context
			ctx := context.WithValue(r.Context(), userCtxKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserFromContext(ctx context.Context) (*db.User, bool) {
	user, ok := ctx.Value(userCtxKey).(db.User)
	return &user, ok
}
