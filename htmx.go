package htmx

import (
	"net/http"
	"strings"
)

type (
	HTMX struct{}
)

func New() *HTMX {
	return &HTMX{}
}

func (s *HTMX) NewHandler(w http.ResponseWriter, r *http.Request) *Handler {
	return &Handler{
		w:       w,
		r:       r,
		request: s.HxHeader(r.Context()),
		response: &HxResponseHeader{
			Headers: make(map[HxResponseKey]string),
		},
	}
}

func HxStrToBool(str string) bool {
	if strings.EqualFold(str, "true") {
		return true
	}

	return false
}

func HxBoolToStr(b bool) string {
	if b {
		return "true"
	}

	return "false"
}
