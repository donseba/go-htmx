package htmx

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	// ContextRequestHeader is the context key for the htmx request header.
	ContextRequestHeader = "htmx-request-header"

	HxRequestHeaderBoosted               HxRequestHeaderKey = "HX-Boosted"
	HxRequestHeaderCurrentURL            HxRequestHeaderKey = "HX-Current-URL"
	HxRequestHeaderHistoryRestoreRequest HxRequestHeaderKey = "HX-History-Restore-Request"
	HxRequestHeaderPrompt                HxRequestHeaderKey = "HX-Prompt" // htmx 4 sends it only with the hx-prompt extension
	HxRequestHeaderRequest               HxRequestHeaderKey = "HX-Request"
	HxRequestHeaderRequestType           HxRequestHeaderKey = "HX-Request-Type" // htmx 4: "full" or "partial"
	HxRequestHeaderSource                HxRequestHeaderKey = "HX-Source"       // htmx 4: the element that triggered the request, tagName#id
	HxRequestHeaderTarget                HxRequestHeaderKey = "HX-Target"       // htmx 2: the id of the target, htmx 4: tagName#id
	HxRequestHeaderTriggerName           HxRequestHeaderKey = "HX-Trigger-Name" // htmx 2 only
	HxRequestHeaderTrigger               HxRequestHeaderKey = "HX-Trigger"      // htmx 2 only, htmx 4 sends HX-Source

	// HxRequestTypeFull is the HX-Request-Type of an htmx 4 request for the whole page,
	// as htmx 4 asks for body swaps, hx-select and history restoration.
	HxRequestTypeFull = "full"

	// HxRequestTypePartial is the HX-Request-Type of an htmx 4 request for a part of the page.
	HxRequestTypePartial = "partial"
)

type (
	HxRequestHeaderKey string

	HxRequestHeader struct {
		HxBoosted               bool
		HxCurrentURL            string
		HxHistoryRestoreRequest bool
		HxPrompt                string
		HxRequest               bool
		HxRequestType           string
		HxSource                string
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
		HxRequestType:           r.Header.Get(HxRequestHeaderRequestType.String()),
		HxSource:                r.Header.Get(HxRequestHeaderSource.String()),
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

// elementID returns the id in an element identifier of htmx 4, tagName#encodeURI(id), decoded.
// It is empty for an element without an id, and the id as sent when it does not decode.
func elementID(identifier string) string {
	_, id, ok := strings.Cut(identifier, "#")
	if !ok {
		return ""
	}

	decoded, err := url.PathUnescape(id)
	if err != nil {
		return id
	}

	return decoded
}
