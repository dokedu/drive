package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// CreateStack takes a list of MiddlewareFunc and returns a single MiddlewareFunc.
func CreateStack(middlewares ...mux.MiddlewareFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
