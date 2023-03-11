package middleware

import (
	"context"
	"net/http"

	"github.com/donseba/go-htmx"
)

func MiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hxh := htmx.HxHeaderRequest{
			HxBoosted:               htmx.HxStrToBool(r.Header.Get("HX-Boosted")),
			HxCurrentURL:            r.Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: htmx.HxStrToBool(r.Header.Get("HX-History-Restore-Request")),
			HxPrompt:                r.Header.Get("HX-Prompt"),
			HxRequest:               htmx.HxStrToBool(r.Header.Get("HX-Request")),
			HxTarget:                r.Header.Get("HX-Target"),
			HxTriggerName:           r.Header.Get("HX-Trigger-Name"),
			HxTrigger:               r.Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
