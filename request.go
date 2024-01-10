package htmx

import (
	"net/http"
)

const ContextRequestHeader = "htmx-request-header"

type (
	HxRequestHeader struct {
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

func HxRequestHeaderFromRequest(r *http.Request) HxRequestHeader {
	return HxRequestHeader{
		HxBoosted:               HxStrToBool(r.Header.Get("HX-Boosted")),
		HxCurrentURL:            r.Header.Get("HX-Current-URL"),
		HxHistoryRestoreRequest: HxStrToBool(r.Header.Get("HX-History-Restore-Request")),
		HxPrompt:                r.Header.Get("HX-Prompt"),
		HxRequest:               HxStrToBool(r.Header.Get("HX-Request")),
		HxTarget:                r.Header.Get("HX-Target"),
		HxTriggerName:           r.Header.Get("HX-Trigger-Name"),
		HxTrigger:               r.Header.Get("HX-Trigger"),
	}
}

func (h *HTMX) HxHeader(r *http.Request) HxRequestHeader {
	header := r.Context().Value(ContextRequestHeader)

	if val, ok := header.(HxRequestHeader); ok {
		return val
	}

	// if the header is not found from the middleware, try and populate it from the request
	return HxRequestHeaderFromRequest(r)
}
