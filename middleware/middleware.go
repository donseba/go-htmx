package middleware

import (
	"context"
	"net/http"

	"github.com/donseba/go-htmx"
)

// MiddleWare is a middleware that adds the htmx request header to the context
// deprecated: htmx will retrieve the headers from the request by itself using htmx.NewHandler(w, r)
func MiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hxh := htmx.HxRequestHeaderFromRequest(r)

		//nolint:staticcheck
		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
