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
		// read request header Authorization
		authHeader := r.Header.Get("Authorization")

		// write auth header to ctx
		ctx := r.Context()
		ctx = context.WithValue(ctx, authKey, authHeader)

		// pass request to next handler
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
