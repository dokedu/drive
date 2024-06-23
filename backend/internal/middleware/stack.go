package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

// CreateStack takes a list of MiddlewareFuncs and returns a single MiddlewareFunc.
func CreateStack(middlewares ...mux.MiddlewareFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
