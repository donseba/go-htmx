package htmx

import (
	"encoding/json"
	"net/http"
)

type (
	Handler struct {
		w           http.ResponseWriter
		r           *http.Request
		request     HxRequestHeader
		response    *HxResponseHeader
		wroteHeader bool
		statusCode  int
	}
)

// Write writes the data to the connection as part of an HTTP reply.
func (h *Handler) Write(data []byte) (n int, err error) {
	for k, v := range h.response.Headers {
		h.w.Header().Set(k.String(), v)
	}

	h.w.WriteHeader(h.statusCode)
	return h.w.Write(data)
}

// WriteHeader sets the response status
func (h *Handler) WriteHeader(code int) {
	if h.wroteHeader {
		return
	}

	h.wroteHeader = true
	h.statusCode = code
}

// Header returns the header map that will be sent by WriteHeader
func (h *Handler) Header() http.Header {
	return h.w.Header()
}

type LocationInput struct {
	Source  string                 `json:"source"`  // source - the source element of the request
	Event   string                 `json:"event"`   //event - an event that "triggered" the request
	Handler string                 `json:"handler"` //handler - a callback that will handle the response HTML
	Target  string                 `json:"target"`  //target - the target to swap the response into
	Swap    string                 `json:"swap"`    //swap - how the response will be swapped in relative to the target
	Values  map[string]interface{} `json:"values"`  //values - values to submit with the request
	Header  map[string]interface{} `json:"headers"` //headers - headers to submit with the request

}

// Location can be used to trigger a client side redirection without reloading the whole page
// https://htmx.org/headers/hx-location/
func (h *Handler) Location(li *LocationInput) error {
	payload, err := json.Marshal(li)
	if err != nil {
		return err
	}

	h.response.Set(HXLocation, string(payload))
	return nil
}

// PushURL pushes a new url into the history stack.
// https://htmx.org/headers/hx-push-url/
func (h *Handler) PushURL(val string) {
	h.response.Set(HXPushUrl, val)
}

// Redirect can be used to do a client-side redirect to a new location
func (h *Handler) Redirect(val string) {
	h.response.Set(HXRedirect, val)
}

// Refresh if set to true the client side will do a full refresh of the page
func (h *Handler) Refresh(val bool) {
	h.response.Set(HXRefresh, HxBoolToStr(val))
}

// ReplaceURL allows you to replace the current URL in the browser location history.
// https://htmx.org/headers/hx-replace-url/
func (h *Handler) ReplaceURL(val string) {
	h.response.Set(HXReplaceUrl, val)
}

// ReSwap allows you to specify how the response will be swapped. See hx-swap for possible values
// https://htmx.org/attributes/hx-swap/
func (h *Handler) ReSwap(val string) {
	h.response.Set(HXReswap, val)
}

// ReTarget a CSS selector that updates the target of the content update to a different element on the page
func (h *Handler) ReTarget(val string) {
	h.response.Set(HXRetarget, val)
}

// Trigger triggers events as soon as the response is received.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) Trigger(val string) {
	h.response.Set(HXTrigger, val)
}

// TriggerAfterSettle trigger events after the settling step.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) TriggerAfterSettle(val string) {
	h.response.Set(HXTriggerAfterSettle, val)
}

// TriggerAfterSwap trigger events after the swap step.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) TriggerAfterSwap(val string) {
	h.response.Set(HXTriggerAfterSwap, val)
}

// Request returns the HxHeaders from the request
func (h *Handler) Request() HxRequestHeader {
	return h.request
}
