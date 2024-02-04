package htmx

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type (
	Handler struct {
		log      Logger
		w        http.ResponseWriter
		r        *http.Request
		request  HxRequestHeader
		response *HxResponseHeader
	}
)

const (
	// StatusStopPolling is the status code that will stop htmx from polling
	StatusStopPolling = 286
)

// IsHxRequest returns true if the request is a htmx request.
func (h *Handler) IsHxRequest() bool {
	return h.request.HxRequest
}

// IsHxBoosted returns true if the request is a htmx request and the request is boosted
func (h *Handler) IsHxBoosted() bool {
	return h.request.HxBoosted
}

// IsHxHistoryRestoreRequest returns true if the request is a htmx request and the request is a history restore request
func (h *Handler) IsHxHistoryRestoreRequest() bool {
	return h.request.HxHistoryRestoreRequest
}

// RenderPartial returns true if the request is an HTMX request that is either boosted or a standard request,
// provided it is not a history restore request.
func (h *Handler) RenderPartial() bool {
	return (h.request.HxRequest || h.request.HxBoosted) && !h.request.HxHistoryRestoreRequest
}

// Write writes the data to the connection as part of an HTTP reply.
func (h *Handler) Write(data []byte) (n int, err error) {
	return h.w.Write(data)
}

// WriteHTML is a helper that writes HTML data to the connection.
func (h *Handler) WriteHTML(html template.HTML) (n int, err error) {
	return h.Write([]byte(html))
}

// WriteString is a helper that writes string data to the connection.
func (h *Handler) WriteString(s string) (n int, err error) {
	return h.Write([]byte(s))
}

// WriteJSON is a helper that writes json data to the connection.
func (h *Handler) WriteJSON(data any) (n int, err error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	return h.Write(payload)
}

// JustWrite writes the data to the connection as part of an HTTP reply.
func (h *Handler) JustWrite(data []byte) {
	_, err := h.Write(data)
	if err != nil {
		h.log.Warn(err.Error())
	}
}

// JustWriteHTML is a helper that writes HTML data to the connection.
func (h *Handler) JustWriteHTML(html template.HTML) {
	_, err := h.WriteHTML(html)
	if err != nil {
		h.log.Warn(err.Error())
	}
}

// JustWriteString is a helper that writes string data to the connection.
func (h *Handler) JustWriteString(s string) {
	_, err := h.WriteString(s)
	if err != nil {
		h.log.Warn(err.Error())
	}
}

// JustWriteJSON is a helper that writes json data to the connection.
func (h *Handler) JustWriteJSON(data any) {
	_, err := h.WriteJSON(data)
	if err != nil {
		h.log.Warn(err.Error())
	}
}

// MustWrite writes the data to the connection as part of an HTTP reply.
func (h *Handler) MustWrite(data []byte) {
	_, err := h.Write(data)
	if err != nil {
		panic(err)
	}
}

// MustWriteHTML is a helper that writes HTML data to the connection.
func (h *Handler) MustWriteHTML(html template.HTML) {
	_, err := h.WriteHTML(html)
	if err != nil {
		panic(err)
	}
}

// MustWriteString is a helper that writes string data to the connection.
func (h *Handler) MustWriteString(s string) {
	_, err := h.WriteString(s)
	if err != nil {
		panic(err)
	}
}

// MustWriteJSON is a helper that writes json data to the connection.
func (h *Handler) MustWriteJSON(data any) {
	_, err := h.WriteJSON(data)
	if err != nil {
		panic(err)
	}
}

// WriteHeader sets the HTTP response header with the provided status code.
func (h *Handler) WriteHeader(code int) {
	h.w.WriteHeader(code)
}

// StopPolling sets the response status to 286, which will stop htmx from polling
func (h *Handler) StopPolling() {
	h.WriteHeader(StatusStopPolling)
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

// ReSwapWithObject allows you to specify how the response will be swapped. See hx-swap for possible values
// https://htmx.org/attributes/hx-swap/
func (h *Handler) ReSwapWithObject(s *Swap) {
	h.ReSwap(s.String())
}

// ReTarget a CSS selector that updates the target of the content update to a different element on the page
func (h *Handler) ReTarget(val string) {
	h.response.Set(HXRetarget, val)
}

// ReSelect a CSS selector that allows you to choose which part of the response is used to be swapped in. Overrides an existing hx-select on the triggering element
func (h *Handler) ReSelect(val string) {
	h.response.Set(HXReselect, val)
}

// Trigger triggers events as soon as the response is received.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) Trigger(val string) {
	h.response.Set(HXTrigger, val)
}

// TriggerWithObject triggers events as soon as the response is received.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) TriggerWithObject(t *Trigger) {
	h.Trigger(t.String())
}

// TriggerAfterSettle trigger events after the settling step.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) TriggerAfterSettle(val string) {
	h.response.Set(HXTriggerAfterSettle, val)
}

// TriggerAfterSettleWithObject trigger events after the settling step.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) TriggerAfterSettleWithObject(t *Trigger) {
	h.TriggerAfterSettle(t.String())
}

// TriggerAfterSwap trigger events after the swap step.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) TriggerAfterSwap(val string) {
	h.response.Set(HXTriggerAfterSwap, val)
}

// TriggerAfterSwapWithObject trigger events after the swap step.
// https://htmx.org/headers/hx-trigger/
func (h *Handler) TriggerAfterSwapWithObject(t *Trigger) {
	h.TriggerAfterSwap(t.String())
}

// Request returns the HxHeaders from the request
func (h *Handler) Request() HxRequestHeader {
	return h.request
}

// ResponseHeader returns the value of the response header
func (h *Handler) ResponseHeader(header HxResponseKey) string {
	return h.response.Get(header)
}
