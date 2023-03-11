package htmx

import (
	"net/http"
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
