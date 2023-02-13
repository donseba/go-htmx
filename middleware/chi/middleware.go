package htmx_chi

import (
	"context"
	"net/http"
	"strings"

	go_htmx "github.com/donseba/go-htmx"
)

func MiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hxh := go_htmx.HxHeaderRequest{
			HxBoosted:               parseStrAsBool(r.Header.Get("HX-Boosted")),
			HxCurrentURL:            r.Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: parseStrAsBool(r.Header.Get("HX-History-Restore-Request")),
			HxPrompt:                r.Header.Get("HX-Prompt"),
			HxRequest:               parseStrAsBool(r.Header.Get("HX-Request")),
			HxTarget:                r.Header.Get("HX-Target"),
			HxTriggerName:           r.Header.Get("HX-Trigger-Name"),
			HxTrigger:               r.Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, go_htmx.ContextRequestHeader, hxh)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func parseStrAsBool(str string) bool {
	if strings.EqualFold(str, "true") {
		return true
	}

	return false
}
