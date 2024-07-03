package middleware

import (
	"context"
	"example/internal/database"
	"example/internal/database/db"
	"net/http"
)

const authKey = "auth"

func Authentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Check if the request has an Authorization header
		authHeader := r.Header.Get("Authorization")

		// Write auth header to request context
		ctx = context.WithValue(ctx, authKey, authHeader)

		// Pass request to next handler
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context, db *database.DB) (*db.User, bool) {
	authHeader := ctx.Value(authKey)
	if authHeader == nil {
		return nil, false
	}
	user, err := db.GLOBAL_UserFindBySessionToken(ctx, authHeader.(string))
	if err != nil {
		return nil, false
	}
	return &user, true
}
