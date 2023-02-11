package htmx

import (
	"context"
	"net/http"
	"strings"
)

const hxHeaderCtx = "hx-header-ctx"

type (
	HxHeader struct {
		HxBoosted               bool
		HxCurrentURL            string
		HxHistoryRestoreRequest bool
		HxPrompt                string
		HxRequest               bool
		HxTarget                string
		HxTriggerName           string
		HxTrigger               string
	}
)

func (s *Service) HxHeaderMiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		hxh := HxHeader{
			HxBoosted:               parseStrAsBool(r.Header.Get("HX-Boosted")),
			HxCurrentURL:            r.Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: parseStrAsBool(r.Header.Get("HX-History-Restore-Request")),
			HxPrompt:                r.Header.Get("HX-Prompt"),
			HxRequest:               parseStrAsBool(r.Header.Get("HX-Request")),
			HxTarget:                r.Header.Get("HX-Target"),
			HxTriggerName:           r.Header.Get("HX-Trigger-Name"),
			HxTrigger:               r.Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, hxHeaderCtx, hxh)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (s *Service) HxHeader(ctx context.Context) HxHeader {
	header := ctx.Value(hxHeaderCtx)

	if val, ok := header.(HxHeader); ok {
		return val
	}

	return HxHeader{}
}
func parseStrAsBool(str string) bool {
	if strings.EqualFold(str, "true") {
		return true
	}

	return false
}
