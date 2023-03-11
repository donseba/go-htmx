package htmx

import (
	"context"
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

func (s *HTMX) HxHeader(ctx context.Context) HxRequestHeader {
	header := ctx.Value(ContextRequestHeader)

	if val, ok := header.(HxRequestHeader); ok {
		return val
	}

	return HxRequestHeader{}
}
