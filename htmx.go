// Package htmx offers a streamlined integration with HTMX in Go applications.
// It implements the standard io.Writer interface and includes middleware support, but it is not required.
// Allowing for the effortless incorporation of HTMX features into existing Go applications.
package htmx

import (
	"errors"
	"github.com/donseba/go-htmx/sse"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var (
	DefaultSwapDuration = time.Duration(0 * time.Millisecond)
	DefaultSettleDelay  = time.Duration(20 * time.Millisecond)

	DefaultNotificationKey   = "showMessage"
	DefaultSSEWorkerPoolSize = 5
)

// this is the default sseManager implementation which is created to handle the server-sent events.
var sseManager sse.Manager

type (
	Logger interface {
		Warn(msg string, args ...any)
	}

	HTMX struct {
		log Logger
	}
)

// New returns a new htmx instance.
func New() *HTMX {
	return &HTMX{
		log: slog.Default().WithGroup("htmx"),
	}
}

// SetLog sets the logger for the htmx instance.
func (h *HTMX) SetLog(log Logger) {
	h.log = log
}

// NewHandler returns a new htmx handler.
func (h *HTMX) NewHandler(w http.ResponseWriter, r *http.Request) *Handler {
	return &Handler{
		w:        w,
		r:        r,
		request:  h.HxHeader(r),
		response: h.HxResponseHeader(w.Header()),
		log:      h.log,
	}
}

// NewSSE creates a new sse manager with the specified worker pool size.
func (h *HTMX) NewSSE(workerPoolSize int) error {
	if sseManager != nil {
		return errors.New("sse manager already exists")
	}

	sseManager = sse.NewManager(workerPoolSize)
	return nil
}

// SSEHandler handles the server-sent events. this is a shortcut and is not the preferred way to handle sse.
func (h *HTMX) SSEHandler(w http.ResponseWriter, r *http.Request, cl sse.Listener) {
	if sseManager == nil {
		sseManager = sse.NewManager(DefaultSSEWorkerPoolSize)
	}

	sseManager.Handle(w, r, cl)
}

// SSESend sends a message to all connected clients.
func (h *HTMX) SSESend(message sse.Envelope) {
	if sseManager == nil {
		sseManager = sse.NewManager(DefaultSSEWorkerPoolSize)
	}

	sseManager.Send(message)
}

// IsHxRequest returns true if the request is a htmx request.
func IsHxRequest(r *http.Request) bool {
	return HxStrToBool(r.Header.Get(HxRequestHeaderRequest.String()))
}

// IsHxBoosted returns true if the request is a htmx request and the request is boosted
func IsHxBoosted(r *http.Request) bool {
	return HxStrToBool(r.Header.Get(HxRequestHeaderBoosted.String()))
}

// IsHxHistoryRestoreRequest returns true if the request is a htmx request and the request is a history restore request
func IsHxHistoryRestoreRequest(r *http.Request) bool {
	return HxStrToBool(r.Header.Get(HxRequestHeaderHistoryRestoreRequest.String()))
}

// RenderPartial returns true if the request is an HTMX request that is either boosted or a hx request,
// provided it is not a history restore request.
func RenderPartial(r *http.Request) bool {
	return (IsHxRequest(r) || IsHxBoosted(r)) && !IsHxHistoryRestoreRequest(r)
}

// HxStrToBool converts a string to a boolean value.
func HxStrToBool(str string) bool {
	return strings.EqualFold(str, "true")
}

// HxBoolToStr converts a boolean value to a string.
func HxBoolToStr(b bool) string {
	if b {
		return "true"
	}

	return "false"
}
