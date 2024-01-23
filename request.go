package htmx

import (
	"net/http"
)

const (
	// ContextRequestHeader is the context key for the htmx request header.
	ContextRequestHeader = "htmx-request-header"

	HxRequestHeaderBoosted               HxRequestHeaderKey = "HX-Boosted"
	HxRequestHeaderCurrentURL            HxRequestHeaderKey = "HX-Current-URL"
	HxRequestHeaderHistoryRestoreRequest HxRequestHeaderKey = "HX-History-Restore-Request"
	HxRequestHeaderPrompt                HxRequestHeaderKey = "HX-Prompt"
	HxRequestHeaderRequest               HxRequestHeaderKey = "HX-Request"
	HxRequestHeaderTarget                HxRequestHeaderKey = "HX-Target"
	HxRequestHeaderTriggerName           HxRequestHeaderKey = "HX-Trigger-Name"
	HxRequestHeaderTrigger               HxRequestHeaderKey = "HX-Trigger"
)

type (
	HxRequestHeaderKey string

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
		HxBoosted:               HxStrToBool(r.Header.Get(HxRequestHeaderBoosted.String())),
		HxCurrentURL:            r.Header.Get(HxRequestHeaderCurrentURL.String()),
		HxHistoryRestoreRequest: HxStrToBool(r.Header.Get(HxRequestHeaderHistoryRestoreRequest.String())),
		HxPrompt:                r.Header.Get(HxRequestHeaderPrompt.String()),
		HxRequest:               HxStrToBool(r.Header.Get(HxRequestHeaderRequest.String())),
		HxTarget:                r.Header.Get(HxRequestHeaderTarget.String()),
		HxTriggerName:           r.Header.Get(HxRequestHeaderTriggerName.String()),
		HxTrigger:               r.Header.Get(HxRequestHeaderTrigger.String()),
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

func (x HxRequestHeaderKey) String() string {
	return string(x)
}
